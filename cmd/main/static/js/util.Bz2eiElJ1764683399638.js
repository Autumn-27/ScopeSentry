const o=(o,e)=>new Promise((s,l)=>{e.setCrossOrigin("Anonymous"),e.load(o,o=>{s(o)},o=>{console.log(o.loaded/o.total*100+"% loaded")},o=>{console.error(o),l(o)})});export{o as l};
