import fs from 'fs'


const data = {
    counter: "0"
}
if(!fs.existsSync("db.json")){
    fs.writeFileSync("db.json",JSON.stringify(data))
}

