package request_context

type clientKey struct{}

func (r Request) ClientID() string {
	clientID, ok := r.ctx.Value(clientKey{}).(string)
	if ok {
		return clientID
	}
	return ""
}
