import express from "express";
import fs from 'fs'
const app = express()

const poll = [
    {id: "1", name: "haki", votes: 0},
    {id: "2", name: "kaima", votes: 0},
    {id: "3", name: "emezie", votes: 0},
    {id: "4", name: "reign", votes: 0},
    {id: "5", name: "collins", votes: 0},
    {id: "6", name: "doll", votes: 0},
]

if(!fs.existsSync("voters.json")){
    fs.writeFileSync("voters.json", "[]")
}

app.get("/", (req, res) => {
    res.json({message: "api running"})
})

app.post("/vote", (req, res) => {
    const {id, pollid} = req.body

    if(!id || !pollid){
        return res.status(400).json({message: "bad request"})
    }

    if(pollid > 6 || pollid < 1){
        return res.status(400).json({message: "invalid option"})
    }

    const voters = JSON.parse(fs.readFileSync("voters.json", "utf-8"))

    const voted = voters.find(x => x.id == id)
    if(voted){
        return res.json({message: "You have voted before"})
    }

    const opt = poll.find(x => x.id == pollid)
    opt.votes = opt.votes  + 1

    const voter = {
        id,
        pollid
    }

    voters.push(voter)

    fs.writeFileSync("voters.json", JSON.stringify(voters, null, 2))

    return res.json({message: "You have voted successfully"})
})

app.get("/view", (req, res) => {
    let data = ' ' 
    poll.forEach(((po) =>{
       data = data + `Contestant ${po.name} has ${po.votes}
       `
    }))
    return res.json({data})
})

app.listen(3000, () => {
    console.log("api running")
})