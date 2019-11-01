package proxy

type Response struct {
	CorrelationId string
	Response []byte
}

func NewResponse(request Request, response []byte) Response {
	return Response{
		CorrelationId: request.CorrelationId,
		Response: response,
	}
}

func (r *Response) Publish(w Worker) {
	w.Responder <- *r
}
