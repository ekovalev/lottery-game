package handler

import (
	"encoding/binary"
	"encoding/json"
	"io/ioutil"
	"net/http"
	_ "strings"

	"go-server/accumulate"
	"go-server/bonus"
	"go-server/stack"

	"github.com/gin-gonic/gin"
	log "github.com/sirupsen/logrus"
	"gopkg.in/go-playground/validator.v9"
)

type placeBetRequest struct {
	ID     string  `json:"id" validate:"required"`
	Amount float64 `json:"amount" validate:"omitempty,gte=0"`
	Bet    []byte  `json:"bet" validate:"required,len=2,dive,gte=0,lte=255"`
}

type betData struct {
	id      string
	gameFee float64
	bet     uint16
}

type Handler struct {
	stack   *stack.BytePairStack
	jackpot *accumulate.JackpotType
	bonus   *bonus.BonusRegistry
}

// New creates handler
func New(st *stack.BytePairStack, jp *accumulate.JackpotType, br *bonus.BonusRegistry) *Handler {
	return &Handler{
		stack:   st,
		jackpot: jp,
		bonus:   br,
	}
}

func (h *Handler) ParseRequestMiddleware(c *gin.Context) {
	body, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		log.WithError(err).Error("Read request body error")
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	request := new(placeBetRequest)
	if err := json.Unmarshal(body, &request); err != nil {
		log.WithFields(log.Fields{"body": body, "as a string": string(body)}).Error("[Handler::ParseRequestMiddleware] Invalid request format")
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		c.Abort()
		return
	}

	// Validating request against constraints
	v9 := validator.New()
	if err := v9.Struct(request); err != nil {
		log.WithField("err", err).Error("[Handler::ParseRequestMiddleware] Request validation failed")
		c.JSON(http.StatusBadRequest, gin.H{
			"status": err.Error(),
		})
		c.Abort()
		return
	}

	userID := request.ID
	amount := request.Amount
	if amount <= 0.0 {
		// Player must have a free bonus game - check it!
		if err := h.bonus.Use(userID); err != nil {
			c.JSON(http.StatusForbidden, gin.H{
				"status": err.Error(),
			})
			c.Abort()
			return
		}
	}

	bet := binary.LittleEndian.Uint16([]byte{request.Bet[0], request.Bet[1]})

	betData := betData{
		id:      userID,
		gameFee: amount,
		bet:     bet,
	}
	log.WithField("betData", betData).Info("[Handler::ParseRequestMiddleware] Writing data to context")

	c.Set("betData", betData)
}

func (h *Handler) Play(c *gin.Context) {
	bet := c.Value("betData").(betData)
	log.WithField("betData", bet).Info("[Handler::Play] Read data from context")

	log.WithField("luckyPairs", h.stack).Info("[Handler::Play] Stack contents")
	log.WithField("jackpot", h.jackpot).Info("[Handler::Play] Jackpot value")
	log.WithField("bonusGames", h.bonus).Info("[Handler::Play] bonusGames contents")

	// We don't require a transaction-like behaviour when the state of the game (stack, jackpot and assigned bonus games)
	// would be "frozen" for a user starting playing; contribuitng to jackpot and drawing a lucky pair of bytes are independent

	// Player has to pay for the game first
	h.jackpot.Add(bet.gameFee)

	// Pop value from stack of luckyPairs
	pair, err := h.stack.Pop()
	if err != nil {
		// Stack is empty - failed to draw a pair
		// Grant the played a free bonus game and abort
		h.bonus.Add(bet.id)
		c.JSON(http.StatusInternalServerError, gin.H{
			"status":  "ERR",
			"error":   err.Error(),
			"message": "Bonus game assigned to the player; try again with empty game fee field",
		})
		return
	}

	if bet.bet == pair {
		// Won! Trying to redeem the jackpot
		reward := h.jackpot.Redeem()

		// Check if we are unlucky and someone else had drained the jackpot just before us
		// For testing purpose consider replaceing the line below with something like to render this event easily observable
		// if reward < 1.0 {
		if reward == 0.0 {
			// Granting free bonus game
			h.bonus.Add(bet.id)
			c.JSON(http.StatusOK, gin.H{
				"status":  "BONUS_GAME",
				"message": "Jackpot empty; bonus game assigned to the player; try again with empty game fee field",
			})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":  "WIN",
			"jackpot": reward,
			"message": "Well done! Congrats on your win!!!",
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"status":  "NO_WIN",
		"message": "Unlucky. Try again next time",
	})
	return
}
