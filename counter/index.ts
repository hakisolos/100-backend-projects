import { Hono } from "hono";
import fs from 'fs'
const app = new Hono()


app.patch("/increment", (c) => {
    let counterDB = JSON.parse(fs.readFileSync("db.json",'utf-8'))
    counterDB.counter++
    fs.writeFileSync("db.json", JSON.stringify(counterDB))
    return c.text('done')
})


app.patch("/decrement", (c) => {
    let counterDB = JSON.parse(fs.readFileSync("db.json",'utf-8'))
    if(counterDB.counter == 0){
        return c.text('decrease, value cant be less than zero')
    }
    counterDB.counter--
    fs.writeFileSync("db.json", JSON.stringify(counterDB))
    return c.text('done')
})

app.get("/pageviews", (c) => {
    let counterDB = JSON.parse(fs.readFileSync("db.json",'utf-8'))
    return c.json({counter: counterDB.counter})
})

app.patch("/reset", (c) => {
    let counterDB = JSON.parse(fs.readFileSync("db.json",'utf-8'))
    counterDB.counter = 0;
    fs.writeFileSync("db.json", JSON.stringify(counterDB))
    return c.text('done')
})

Bun.serve({
    fetch: app.fetch,
    port: 3000
})

console.log('app running')