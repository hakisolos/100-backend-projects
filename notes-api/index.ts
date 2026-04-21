import { Hono } from "hono";
import { ExistNote, GenerateRandomId, InitDB, ValidOrders } from "./utils";
import fs from 'fs'
import type { Note, UpdateNoteRequest } from "./utils/types";
const app = new Hono()

InitDB()

app.get("/", (c) => {
    return c.json({message: "api running"},200)
})

app.post("/notes/add", async(c) => {
    const { title,  body } = await c.req.json()
    if(!title || !body){
        return c.json({error: "bad request"},400)
    }
    const notes: Array<Note> = JSON.parse(fs.readFileSync("notes.json", "utf-8"))
    const data: Note = {
        id: GenerateRandomId(15),
        title,
        body,
        created_at: Date.now(),
        updated_at: Date.now()
    }
    notes.push(data)
    fs.writeFileSync("notes.json",JSON.stringify(notes,null,2))
    return c.json({message: "note added successfullt"},200)
})

app.delete("/notes/delete", (c) => {
    const id = c.req.query("id")
    if(!id){
        return c.text('id needed',400)
    }
    let notes: Array<Note> = JSON.parse(fs.readFileSync("notes.json", "utf-8"))
    const exist = ExistNote(id)
    if(!exist){
        return c.json({"error": "Note does not exist"})
    }
    notes = notes.filter(fek => fek.id !== id)
    fs.writeFileSync("notes.json",JSON.stringify(notes,null,2))
    return c.json({message: "note deleted successfullt"},200)
})

app.get("/notes/getOne", (c) => {
    const id = c.req.query("id")
    if(!id){
        return c.text('id needed',400)
    }
    let notes: Array<Note> = JSON.parse(fs.readFileSync("notes.json", "utf-8"))
    const exist = ExistNote(id)
    if(!exist){
        return c.json({"error": "Note does not exist"})
    }
    const data = notes.find(fek => fek.id == id)
    return c.json({message: "note found", data})
})

app.get("/notes/get", c => {
    const order = c.req.query("order")
    let notes: Array<Note> = JSON.parse(fs.readFileSync("notes.json", "utf-8"))
  
    if(notes.length == 0){
        return c.text('no notes yet')
    }
    if(order){
        if(!ValidOrders.includes(order)){
            return c.text(`invalid parameter, use "newest" or "oldest"`)
        }
        if(order == "newest"){
            notes.sort((ha,ki) => ki.updated_at - ha.updated_at)
            let data = notes
            return c.json({data})
        }
        if(order == "oldest"){
            notes.sort((ha,ki) => ha.updated_at - ki.updated_at)
            let data = notes
            return c.json({data})
        }
    }
    notes.sort((ha,ki) => ki.updated_at - ha.updated_at)
    let data = notes
    return c.json({data})

})

app.patch("/notes/update", async(c) => {
    const body: UpdateNoteRequest = await c.req.json()
    let notes: Array<Note> = JSON.parse(fs.readFileSync("notes.json", "utf-8"))
    if(!(ExistNote(body.id))){
        return c.text('note does not exist',404)
    }
    if(!body.id){
        return c.text('id required',400)
    }
    if(body.title && body.body){
        let note = notes.find(n => n.id == body.id)
        if(!note){
            return c.text('note does not exist',404)
        }
        note.title = body.title
        note.body = body.body
        note.updated_at = Date.now()
        fs.writeFileSync("notes.json",JSON.stringify(notes,null,2))
        return c.json({message: "note updated successfully"},200)
    } 
    if(body.body){
        let note = notes.find(n => n.id == body.id)
        if(!note){
            return c.text('note does not exist',404)
        }
        note.body = body.body
        note.updated_at = Date.now()
        fs.writeFileSync("notes.json",JSON.stringify(notes,null,2))
        return c.json({message: "note updated successfully"},200)
    }
    if(body.title){
        let note = notes.find(n => n.id == body.id)
        if(!note){
            return c.text('note does not exist',404)
        }
        note.title = body.title
        note.updated_at = Date.now()
        fs.writeFileSync("notes.json",JSON.stringify(notes,null,2))
        return c.json({message: "note updated successfully"},200)
    }
})

Bun.serve({
    fetch: app.fetch,
    port: 3000
})
console.log("app running")