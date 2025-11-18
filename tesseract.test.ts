// load jpn-test.png as base64 and send to server

import * as fs from 'fs';

const imageBuffer = fs.readFileSync('jpn-test.png');
const base64 = imageBuffer.toString('base64');

fetch("http://localhost:3000/upload", {
  method: "POST",
  headers: {
    "Content-Type": "application/json"
  },
  body: JSON.stringify({
    base64Image: `data:image/png;base64,${base64}`,
    lang: "jpn"
  })
}).then(res => res.json()).then(data => {
  fs.writeFileSync('tesseract_example_output.json', JSON.stringify(data, null, 2), 'utf8');
  console.log('Saved response to output.json');
});