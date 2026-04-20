import e from "express";
import fs from 'fs'
import { ExistTask, Statuses } from "./utils.js";

const app = e()

app.use(e.json())

if(!fs.existsSync("database.json")){
    fs.writeFileSync("database.json","[]")
}

app.get("/",(req, res) => {
    return res.json({message: "api running"})
})

app.post("/create", (req, res) => {
    const {title,description} = req.body
    if (!title || !description){
        return res.status(400).json({error: "bad request"})
    }

    let list = JSON.parse(fs.readFileSync("database.json",'utf-8'))
    
    list.push({
        id: list.length + 1,
        title,
        description,
        status: "incomplete"
    })

    fs.writeFileSync("database.json",JSON.stringify(list,null,2))
    return res.json({message: "todo added"})
})

app.get("/task/get", (req, res) => {
    const id = req.query.id
    if(!id){
        let tasks = JSON.parse(fs.readFileSync("database.json",'utf-8'))
        return res.status(400).json({data: tasks})
    }
    const fek = ExistTask(id)
    if(!fek.data){
        return res.status(404).json({error: `task of id: ${id} not found`})
    } else {
        return res.json({data: fek.data})
    }
})


app.get("/task/delete", (req, res) => {
    const id = req.query.id
    if(!id){
        return res.status(400).json({error: "bad request"})
    }
    const fek = ExistTask(id)
    let tasks = fek.tasks
    let data = fek.data
    if(!data){
        return res.status(404).json({error: `task of id: ${id} not found`})
    } else {
        tasks = tasks.filter(feks => feks == id)
        return res.status(200).json({message: `task of id: ${id} deleted successfully`})
    }
})


app.get("/task/list", (req, res) => {
    let tasks = JSON.parse(fs.readFileSync("database.json",'utf-8'))
    if(tasks.length == 0){
        return res.status(404).json({error: "no tasks found"})
    }
    const stat = req.query.status
    if(stat){
        if(!Statuses.includes(stat)){
            console.log(stat)
            console.log(Statuses.includes(stat))
            console.log(Statuses)
            return res.status(400).json({error: "bad request"})
        }
        let data = []
        for(const task of tasks){
            if(task.status == stat){
                data.push(task)
            }
        }
        if (!data){
            return res.status(400).json({error: `no ${stat} tasks`})
        }
        return res.status(200).json({data})
    }
    return res.status(200).json({data: tasks})
})


app.patch("/tasks/update/complete", (req, res) => {
    const id = req.query.id
    if(!id){
        return res.status(400).json({error: "bad request"})
    }
    const fek = ExistTask(id)
    let tasks = fek.tasks
    let data = fek.data
    if(!data){
        return res.status(404).json({error: `task of id: ${id} not found`})
    }else{
        if(data.status == "complete"){
            return res.json({error: "task already complete"})
        } else {
            data.status = "complete"
        }
        tasks = tasks.filter(a => a == data.id)
        tasks.push(data)
        fs.writeFileSync("database.json",JSON.stringify(tasks,null,2))
    }
    return res.status(200).json({message: "updated successfully"})
})


app.patch("/tasks/update/incomplete", (req, res) => {
    const id = req.query.id
    if(!id){
        return res.status(400).json({error: "bad request"})
    }
    const fek = ExistTask(id)
    let tasks = fek.tasks
    let data = fek.data
    if(!data){
        return res.status(404).json({error: `task of id: ${id} not found`})
    }else{
        if(data.status == "incomplete"){
            return res.json({error: "task already incomplete"})
        } else {
            data.status = "incomplete"
        }
        tasks = tasks.filter(a => a == data.id)
        tasks.push(data)
        fs.writeFileSync("database.json",JSON.stringify(tasks,null,2))
    }
    return res.status(200).json({message: "updated successfully"})
})



app.listen(3000)