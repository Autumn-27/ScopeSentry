import{_ as e}from"./ContentWrap.vue_vue_type_script_setup_true_lang-BgCcO80c.js";import{d as t,z as a,a5 as i,r as s,o,i as r,aG as l,a as n,j as m,dr as p,e as d,dh as u,l as c,C as h,w as b,y as g,t as w}from"./index-B4Nyjh3H.js";import{E as f}from"./el-image-viewer-CLYRUyY9.js";import"./el-card-wk8kGoMw.js";import"./el-tooltip-l0sNRNKZ.js";import"./el-popper-DBt0ZFPB.js";import"./debounce-4WtFbZlf.js";const _=t({__name:"ImageViewer",props:{urlList:{type:Array,default:()=>[]},zIndex:a.number.def(200),initialIndex:a.number.def(0),infinite:a.bool.def(!0),hideOnClickModal:a.bool.def(!1),teleported:a.bool.def(!1),show:a.bool.def(!1)},setup(e){const t=e,a=i((()=>{const e={...t};return delete e.show,e})),p=s(t.show),d=()=>{p.value=!1};return(e,t)=>p.value?(o(),r(n(f),l({key:0},a.value,{onClose:d}),null,16)):m("",!0)}});let x=null;const j=t({__name:"ImageViewer",setup(t){const{t:a}=c(),i=()=>{!function(e){if(!p)return;const{urlList:t,initialIndex:a=0,infinite:i=!0,hideOnClickModal:s=!1,teleported:o=!1,zIndex:r=2e3,show:l=!0}=e,n={},m=document.createElement("div");n.urlList=t,n.initialIndex=a,n.infinite=i,n.hideOnClickModal=s,n.teleported=o,n.zIndex=r,n.show=l,document.body.appendChild(m),x=d(_,n),u(x,m)}({urlList:["https://images6.alphacoders.com/657/thumbbig-657194.webp","https://images3.alphacoders.com/677/thumbbig-677688.webp","https://images4.alphacoders.com/200/thumbbig-200966.webp","https://images5.alphacoders.com/657/thumbbig-657248.webp","https://images3.alphacoders.com/679/thumbbig-679917.webp","https://images3.alphacoders.com/737/thumbbig-73785.webp"]})};return(t,s)=>{const l=h("BaseButton");return o(),r(n(e),{title:n(a)("imageViewerDemo.imageViewer"),message:n(a)("imageViewerDemo.imageViewerDes")},{default:b((()=>[d(l,{type:"primary",onClick:i},{default:b((()=>[g(w(n(a)("imageViewerDemo.open")),1)])),_:1})])),_:1},8,["title","message"])}}});export{j as default};