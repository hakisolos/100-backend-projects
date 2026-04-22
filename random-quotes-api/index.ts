import {Hono} from 'hono'
import fs from 'fs'
import type { QuotesDB } from './utils/types'
import { randomItem } from './utils'
const app = new Hono()


app.get("/quotes/random", (c) => {
  const quotes: Array<QuotesDB> = JSON.parse(fs.readFileSync("quotes.json", 'utf-8'))
  const category = c.req.query("category")
  const author = c.req.query("author")

  if(category && author){
    const sorted = []
    for(let quote of quotes){
      if(quote.category == category && quote.author == author){
        sorted.push(quote)
      }
    }
    if(sorted.length == 0){
      return c.text(`no result found for author: ${author} and category: ${category}`)
    } else{
      let data = randomItem(sorted)
      return c.json({data})
    }
  }
  if(category){
    const sorted = []
    for(let quote of quotes){
      if(quote.category == category){
        sorted.push(quote)
      }
    }
    if(sorted.length == 0){
      return c.text(`no result found for category: ${category}`)
    } else{
      let data = randomItem(sorted)
      return c.json({data})
    }
  }
  if(author){
    const sorted = []
    for(let quote of quotes){
      if(quote.author == author){
        sorted.push(quote)
      }
    }
    if(sorted.length == 0){
      return c.text(`no result found for author: ${author}`)
    } else{
      let data = randomItem(sorted)
      return c.json({data})
    }
  }
  let data = randomItem(quotes)
  return c.json({data})
  
})

Bun.serve({
  fetch: app.fetch,
  port: 3000
})

console.log('app running')