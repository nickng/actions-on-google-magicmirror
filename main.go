package backend

import (
	"crypto/subtle"
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

func init() {
	http.HandleFunc("/", Index)
	http.HandleFunc("/action", Auth(Action))
}

const indexTmpl = `<!DOCTYPE html>
<html>
<head>
<title>Magic Mirror</title>
<link rel='stylesheet' href='/static/style.css'/>
</head>
<body>
<h1>Magic Mirror</h1>
<div id='mirror'><img src="/static/mirror.jpg"/><div>
<small>source: https://www.flickr.com/photos/harshlight/6988819332</small>
</body>
</html>`

// Index is a static handler for the landing page request.
func Index(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, indexTmpl)
}

// Context of user's request.
//
// https://api.ai/docs/reference/agent/contexts
type Context struct {
	Name       string            `json:"name"`
	Lifespan   int               `json:"lifespan"`
	Parameters map[string]string `json:"parameters",omitempty`
}

// Request is the data structure for the incoming JSON data.
type Request struct {
	Lang   string `json:lang`
	Status struct {
		Code      int    `json:"code"`
		ErrorType string `json:"errorType"`
	} `json:"status"`
	Timestamp time.Time `json:"timestamp"`
	SessionID string    `json:"sessionId"`
	Result    struct {
		Parameters    map[string]string `json:"parameters"`
		Contexts      []Context         `json:"contexts"`
		ResolvedQuery string            `json:"resolvedQuery"`
		Source        string            `json:"source"`
		Speech        string            `json:"speech"`
		Fulfillment   struct {
			Messages []struct {
				Type   int    `json:"type"`
				Speech string `json:"speech"`
			} `json:"messages"`
			Speech string `json:"speech"`
		} `json:"fulfillment"`
		ActionIncomplete bool   `json:"actionIncomplete"`
		Action           string `json:"action"`
		Metadata         struct {
			IntentID                  string `json:"intentId"`
			WebhookForSlotFillingUsed string `json:"webhookForSlotFillingUsed"`
			IntentName                string `json:"intentName"`
			WebhookUsed               string `json:"webhookUsed"`
		} `json:"metadata"`
		Score float64 `json:"score"`
	} `json:"result"`
	ID              string `json:"id"`
	OriginalRequest struct {
		Source string `json:"source"`
		Data   struct {
			Inputs []struct {
				RawInputs []struct {
					Query     string `json:"query"`
					InputType int    `json:"input_type"`
				} `json:"raw_inputs"`
				Intent    string `json:"intent"`
				Arguments []struct {
					TextValue string `json:"text_value"`
					RawText   string `json:"raw_text"`
					Name      string `json:"name"`
				} `json:"arguments"`
			} `json:"inputs"`
			User struct {
				UserID  string `json:"user_id"`
				Profile struct {
					DisplayName string `json:"display_name"`
					GivenName   string `json:"given_name"`
					FamilyName  string `json:"family_name"`
				} `json:"profile"`
				AccessToken string `json:"access_token"`
			} `json:"user"`
			Conversation struct {
				ConversationToken string `json:"conversation_token"`
				ConversationID    string `json:"conversation_id"`
				Type              int    `json:"type"`
			} `json:"conversation"`
			Device struct {
				Location struct {
					Coordinates struct {
						Latitude  float64 `json:"latitude"`
						Longitude float64 `json:"longitude"`
					} `json:"coordinates"`
					FormattedAddress string `json:"formatted_address"`
					ZipCode          string `json:"zip_code"`
					City             string `json:"city"`
				} `json:"location"`
			} `json:"device"`
		} `json:"data"`
	} `json:"originalRequest"`
}

// Response is the data structure for replies to the server.
type Response struct {
	Speech      string `json:"speech"`
	DisplayText string `json:"displayText"`
	Data        struct {
		Google struct {
			ExpectUserResponse bool `json:"expect_user_response"`
			IsSSML             bool `json:"is_ssml"`
			PermissionsRequest struct {
				OptContext  string   `json:"opt_context"`
				Permissions []string `json:"permissions"`
			} `json:"permissions_request"`
		} `json:"google",omitempty`
	} `json:"data"`
	ContextOut []struct {
		Name       string `json:"name"`
		Lifespan   int    `json:"lifespan"`
		Parameters struct {
			City string `json:"city"`
		} `json:"parameters"`
	} `json:"contextOut"`
	Source        string `json:"source"`
	FollowupEvent struct {
		Name string `json:"name"`
		Data map[string]string
	} `json:"followupEvent"`
}

// Action is a handler for incoming action.
func Action(w http.ResponseWriter, r *http.Request) {
	if r.Body == nil {
		http.Error(w, "Please send a request body", 400)
		return
	}
	var request Request
	if err := json.NewDecoder(r.Body).Decode(&request); err != nil {
		http.Error(w, err.Error(), 500)
		return
	}
	var responseText string
	if choice, ok := request.Result.Parameters["mirror-response"]; ok {
		responseText = generateResponse(choice)
	} else {
		responseText = "I don't quite understand"
	}
	res := Response{
		DisplayText: responseText,
		Speech:      responseText,
	}
	json.NewEncoder(w).Encode(res)
}

// Auth wraps a HandlerFunc with required authentication.
func Auth(fn http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, pass, ok := r.BasicAuth()
		if !ok || !validUser(user, pass) {
			http.Error(w, "Unauthorised", http.StatusUnauthorized)
			return
		}
		fn(w, r)
	}
}

const (
	username = "MagicMirror"
	password = "rorriMcigaM"
)

func validUser(user, pass string) bool {
	return subtle.ConstantTimeCompare([]byte(user), []byte(username)) == 1 &&
		subtle.ConstantTimeCompare([]byte(pass), []byte(password)) == 1
}

func generateResponse(param string) string {
	if param == "fair" {
		return "My Queen, you are the fairest in the land"
	}
	return "Snow White is a thousand times more beautiful than you"
}
