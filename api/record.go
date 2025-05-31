package api

import (
	"database/sql"
	"encoding/json"
	"net/http"

	db "alabast/db/sqlc"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

type createRecordRequest struct {
	Username string          `json:"username" binding:"required"`
	Content  json.RawMessage `json:"content" binding:"required"`
}

func (server *Server) createRecord(ctx *gin.Context) {
	var req createRecordRequest
	if err := ctx.ShouldBindJSON(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.CreateRecordParams{
		Username: req.Username,
		Content:  req.Content,
	}

	record, err := server.store.CreateRecord(ctx, arg)
	if err != nil {
		if pqErr, ok := err.(*pq.Error); ok {
			switch pqErr.Code.Name() {
			case "foreign_key_violation", "unique_violation":
				ctx.JSON(http.StatusForbidden, errorResponse(err))
				return
			}
		}
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, record)
}

type getRecordRequest struct {
	ID int64 `uri:"id" binding:"required,min=1"`
}

func (server *Server) getRecord(ctx *gin.Context) {
	var req getRecordRequest
	if err := ctx.ShouldBindUri(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	record, err := server.store.GetRecord(ctx, req.ID)
	if err != nil {
		if err == sql.ErrNoRows {
			ctx.JSON(http.StatusNotFound, errorResponse(err))
			return
		}

		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, record)
}

type listRecordRequest struct {
	PageID   int32 `form:"page_id" binding:"required,min=1"`
	PageSize int32 `form:"page_size" binding:"required,min=5,max=10"`
}

func (server *Server) listRecords(ctx *gin.Context) {
	var req listRecordRequest
	if err := ctx.ShouldBindQuery(&req); err != nil {
		ctx.JSON(http.StatusBadRequest, errorResponse(err))
		return
	}

	arg := db.ListRecordsParams{
		Limit:  req.PageSize,
		Offset: (req.PageID - 1) * req.PageSize,
	}

	records, err := server.store.ListRecords(ctx, arg)
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, errorResponse(err))
		return
	}

	ctx.JSON(http.StatusOK, records)
}
