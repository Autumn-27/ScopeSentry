import{_ as t}from"./ContentWrap.vue_vue_type_script_setup_true_lang-CTGngVvL.js";import{d as e,l as s,r as i,o as a,i as r,w as o,e as l,a as p,f as n,t as m,y as c}from"./index-CLgjJy2L.js";import{_ as d}from"./Table.vue_vue_type_script_lang-DJc-YjoG.js";import{g as j}from"./index-DPyLJGFH.js";import{E as u}from"./el-link-B6uLfNu7.js";import{E as x}from"./el-divider-B2v5rJA4.js";import"./el-card-CdvbZBex.js";import"./el-tooltip-l0sNRNKZ.js";import"./el-popper-CaFR7Cjt.js";import"./el-table-column-DrBbA1UV.js";import"./index-Cu272xPu.js";import"./debounce-CIcVECfT.js";import"./el-checkbox-DbDCqMQu.js";import"./isArrayLikeObject-CRQFYJ_t.js";import"./raf-BLHXPmBe.js";import"./el-tag-D6hIb9V5.js";import"./el-pagination-uR_gQTrf.js";import"./el-select-CFiSjX0y.js";import"./strings-BAPUa2fl.js";import"./useInput-DX4to9HK.js";import"./el-image-viewer-BF14MS9t.js";import"./el-empty-D1z8uu7s.js";import"./tsxHelper-Cn3mpvqc.js";import"./el-dropdown-item-Dcg5gyE9.js";import"./castArray-CcAiHetO.js";import"./refs-Dw-H7609.js";import"./index-D38LGGz3.js";import"./index-1Oi9uy0T.js";const f={class:"flex cursor-pointer"},v={class:"pr-16px"},_=["src"],g={class:"mb-12px font-700 font-size-16px"},y={class:"line-clamp-3 font-size-12px"},b={class:"flex justify-center items-center"},w=["onClick"],k=["onClick"],C=e({__name:"CardTable",setup(e){const{t:C}=s(),h=i(!0);let z=i([]);(async t=>{const e=await j(t||{pageIndex:1,pageSize:10}).catch((()=>{})).finally((()=>{h.value=!1}));e&&(z.value=e.data.list)})();return(e,s)=>(a(),r(p(t),{title:p(C)("tableDemo.cardTable")},{default:o((()=>[l(p(d),{columns:[],data:p(z),loading:h.value,"custom-content":"","card-wrap-style":{width:"200px",marginBottom:"20px",marginRight:"20px"}},{content:o((t=>[n("div",f,[n("div",v,[n("img",{src:t.logo,class:"w-48px h-48px rounded-[50%]",alt:""},null,8,_)]),n("div",null,[n("div",g,m(t.name),1),n("div",y,m(t.desc),1)])])])),"content-footer":o((t=>[n("div",b,[n("div",{class:"flex-1 text-center",onClick:()=>{}},[l(p(u),{underline:!1},{default:o((()=>[c("操作一")])),_:1})],8,w),l(p(x),{direction:"vertical"}),n("div",{class:"flex-1 text-center",onClick:()=>{}},[l(p(u),{underline:!1},{default:o((()=>[c("操作二")])),_:1})],8,k)])])),_:1},8,["data","loading"])])),_:1},8,["title"]))}});export{C as default};