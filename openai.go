package main

import (
    "bytes"
    "html/template"
    "encoding/json"
    "fmt"
    "net/http"
)

type MyRequest struct {
    Prompt string `json:"prompt"`
}
type MyResponse struct {
    Prompt string `json:"prompt"`
    Content string `json:"content"`
}

type OpenAIRequest struct {
    Model string `json:"model"`
    Messages [1]struct {
	Role string `json:"role"`
	Content string `json:"content"`
    } `json:"messages"`
    Temperature float32 `json:"temperature"`
}

type OpenAIResponse struct {
    ID string `json:"id"`
    Object string `json:"object"`
    Created int `json:"created"`
    Model string `json:"model"`
    Usage struct {
	Prompt_tokens int `json:"prompt_tokens"`
	Completion_tokens int `json:"completion_tokens"`
	Total_tokens int `json:"total_tokens"`
    } `json:"usage"`
    Choices [1]struct {
        Message struct {
	    Role string `json:"role"`
	    Content string `json:"content"`
    } `json:"message"`
	Finish_reason string `json:"finish_reason"`
	Index int `json:"index"`
    } `json:"choices"`
}

func main() {
    http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        if r.Method == "GET" {
            t, _ := template.ParseFiles("index.html")
            t.Execute(w, nil)
        } else if r.Method == "POST" {
            //body, _ := ioutil.ReadAll(r.Body)
            //defer r.Body.Close()

	    //fmt.Print("recv data:%s", string(body))
            //var req MyRequest
            //json.Unmarshal(body, &req)

            client := &http.Client{}
	    var aiReq OpenAIRequest 
	    aiReq.Model = "gpt-3.5-turbo"
	    aiReq.Temperature = 0.7

            aiReq.Messages[0].Role = "user"
            aiReq.Messages[0].Content = r.PostFormValue("prompt")
	    data, _ := json.Marshal(aiReq) 

	    fmt.Print("send data:%s", string(data))

            httpReq, _ := http.NewRequest("POST", "https://api.openai.com/v1/chat/completions", bytes.NewBuffer([]byte(data)))
            httpReq.Header.Set("Content-Type", "application/json")
            httpReq.Header.Set("Authorization", "Bearer xxx")

            resp, _ := client.Do(httpReq)
            defer resp.Body.Close()

	    //for Debug
	    //respBody, _ := ioutil.ReadAll(resp.Body)
	    //fmt.Println(string(respBody))

            var res OpenAIResponse
            err := json.NewDecoder(resp.Body).Decode(&res)
	    if err != nil {
	        w.Write([]byte(fmt.Sprintf("decode failed.%s", err.Error())))	
	    }

	    var myResp MyResponse
	    myResp.Prompt = r.PostFormValue("prompt")
	    myResp.Content = res.Choices[0].Message.Content

	    t, _ := template.ParseFiles("result.html")
	    t.Execute(w, myResp)

            //w.Header().Set("Content-Type", "application/json")
            //w.Write([]byte(res.Choices[0].Message.Content))
        }
    })

    http.ListenAndServe(":8080", nil)
}
