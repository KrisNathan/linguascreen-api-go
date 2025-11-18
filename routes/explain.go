package routes

import (
	"linguascreen/models"
	"net/http"

	"github.com/gin-gonic/gin"
)

type ExplainRoutes struct {
	llmService       models.LLMService
	nmtService       models.NMTService
	memoryRepository models.MemoryRepository
}

func NewExplainRoutes(llmService models.LLMService, nmtService models.NMTService, memoryRepository models.MemoryRepository) *ExplainRoutes {
	return &ExplainRoutes{
		llmService:       llmService,
		nmtService:       nmtService,
		memoryRepository: memoryRepository,
	}
}

type PostExplainRequest struct {
	SelectedText string `json:"selected_text" binding:"required" example:"Hello world"`
	TargetLang   string `json:"target_lang,omitempty" example:"es"`
	UserLang     string `json:"user_lang,omitempty" example:"en"`
}

type PostExplainResponse struct {
	Result models.ExplainResult `json:"result"`
}

// Post godoc
// @Summary Explain selected text
// @Description Explain and translate selected text using LLM and NMT
// @Accept json
// @Produce json
// @Param request body PostExplainRequest true "Explain request"
// @Success 200 {object} PostExplainResponse
// @Failure 400 {object} map[string]string
// @Failure 500 {object} map[string]string
// @Router /explain [post]
func (r *ExplainRoutes) Post(c *gin.Context) {
	var req PostExplainRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	// Rearrange text
	text, rearrangeErr := r.llmService.RearrangeText(req.SelectedText)
	if rearrangeErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": rearrangeErr.Error()})
		return
	}

	// Load memories
	memories, repoErr := r.memoryRepository.GetMemories(0)
	if repoErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": repoErr.Error()})
		return
	}

	// NMT Translation
	nmtTranslation, nmtErr := r.nmtService.TranslateText(text, req.UserLang)
	if nmtErr != nil {
		nmtTranslation = ""
	}

	// Explain
	explanation, explainErr := r.llmService.Explain(text, req.TargetLang, nmtTranslation, req.UserLang, memories)
	if explainErr != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": explainErr.Error()})
		return
	}

	// Respond
	c.IndentedJSON(http.StatusOK, PostExplainResponse{Result: *explanation})
}
