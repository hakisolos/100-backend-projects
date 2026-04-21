import crypto from 'crypto'
import fs from 'fs'
import type { Note } from './types'

export function GenerateRandomId(length: number) : string {
    return crypto.randomBytes(length).toString("hex")
}


export function InitDB() : void {
    if(!fs.existsSync("notes.json")){
        fs.writeFileSync("notes.json","[]")
    }
}

export function ExistNote(id: string): boolean {
    const notes: Array<Note> = JSON.parse(fs.readFileSync("notes.json", "utf-8"))
    var note
    for(let n of notes){
        if(n.id == id){
            note = n
        }
    }
    if(!note){
        return false
    }
    return true
}

export const ValidOrders: Array<string> = ["oldest","newest"]