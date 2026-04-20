import fs from 'fs'
export function ExistTask(id){
    let tasks = JSON.parse(fs.readFileSync("database.json",'utf-8'))
    let fek
    for(let task of tasks){
         if(task.id == id){
                fek = task
        }
    }
    if(!fek){
        return fek
    }
    return {tasks, data: fek}
}

export const Statuses = ["complete","incomplete"]