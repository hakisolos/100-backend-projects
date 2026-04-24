/* 
user passes in url  with titles and tags(optional)
we can fetch all or filter by tag
we can delete bookmark
so we have 
/bookmark/save POST (url: string, title: string, tags?: string[])
/bookmark GET ?tag=
/bookmark DELETE (id)
*/

import { Hono } from "hono";
import { existsSync, readFileSync, writeFileSync } from "node:fs";
import type { getreq, UrlDB } from "./src/types";
import { randomBytes } from "node:crypto";
const app = new Hono()

if(!existsSync("urls.json")){
    writeFileSync("urls.json", "[]")
}

app.post("/bookmark/save", async(c) => {
    let {url, title, tags} = await c.req.json()
    if(!url || !title){
        return c.text("url or title or both missing")
    }
    if(!tags){
        tags = []
    }
    const urlDatabase: Array<UrlDB> = JSON.parse(readFileSync("urls.json","utf-8"))
    urlDatabase.push({
        id: randomBytes(10).toString('hex'),
        url,
        title,
        tags
    })
    writeFileSync("urls.json",JSON.stringify(urlDatabase,null,2))
    return c.text('bookmark saved successfuly')
})

app.get("/bookmark/", (c) => {
    const tag = c.req.query("tag")
    const urlDatabase: Array<UrlDB> = JSON.parse(readFileSync("urls.json","utf-8"))
    let data;
    if(tag){
        data = urlDatabase.filter(something => something.tags?.includes(tag))
    } else {
        data = urlDatabase
    }
    return c.json(data)
})

app.delete("/bookmark/", (c) => {
    const id = c.req.query("id")
    if(!id){
        return c.text("url required")
    }
    let data;
    const urlDatabase: Array<UrlDB> = JSON.parse(readFileSync("urls.json","utf-8"))
    for(let something of urlDatabase){
        if(something.id == id){
            data = something
        }
    } 
    if(!data){
        return c.text("Bookmark Not Found")
    }
    urlDatabase.filter(smth => smth.id == id)
    writeFileSync("urls.json", JSON.stringify(urlDatabase))
    return c.text("Bookmark Deleted Successfully")
})


Bun.serve({
    fetch: app.fetch,
    port: 3000
})
console.log("app running");
