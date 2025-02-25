import{_ as e}from"./ContentWrap.vue_vue_type_script_setup_true_lang-B-a8-WPJ.js";import{d as a,r as l,s as t,e as o,v as i,B as n,C as s,D as r,G as u,F as p,K as c,H as d,N as m,y as g,o as f,c as v,w as y,a as h,z as b,t as _,A as x,f as j,I as F,E as S,j as k,J as w,l as C,M as E}from"./index-C6fb_XFi.js";import{a as A,E as V}from"./el-col-Dl4_4Pn5.js";import{E as D}from"./el-text-BnUG9HvL.js";import{E as U}from"./el-tag-C_oEQYGz.js";import"./el-tooltip-l0sNRNKZ.js";import{E as B}from"./el-popper-CeVwVUf9.js";import{E as I}from"./el-upload-DFauS7op.js";import"./el-progress-sY5OgffI.js";import{E as P,a as z,b as H}from"./el-dropdown-item-DpH7Woj3.js";import{_ as L}from"./Table.vue_vue_type_script_lang-7Pp5E_zy.js";import{u as R}from"./useTable-CijeIiBB.js";import{u as W}from"./useIcon-BxqaCND-.js";import{_ as T}from"./Dialog.vue_vue_type_style_index_0_lang-DjaYHddI.js";import{d as M,g as N,u as O,r as J,a as $,b as G,e as K,c as Q}from"./index-DWlzJn9A.js";import X from"./detail-Dkm4C1dc.js";import"./el-card-B37ahJ8o.js";import"./index-BWEJ0epC.js";import"./el-pagination-FWx5cl5J.js";import"./el-select-vbM8Rxr1.js";import"./strings-BiUeKphX.js";import"./useInput-IB6tFdGu.js";import"./debounce-BwgdhaaK.js";import"./castArray-DRqY4cIf.js";import"./refs-3HtnmaOD.js";import"./el-table-column-C9CkC7I1.js";import"./el-checkbox-CvJzNe2E.js";import"./isArrayLikeObject-zBQ5eq7G.js";import"./raf-DGOAeO92.js";import"./el-image-viewer-DrhdpOg4.js";import"./el-empty-jJjDxScx.js";import"./tsxHelper-CeCzRM_x.js";import"./index-ghAu5K8t.js";import"./index-CnCQNuY4.js";import"./el-form-C2Y6uNCj.js";import"./index-CZoUTVkP.js";const Y={class:"flex flex-wrap gap-3 items-center"},Z={href:"https://plugin.scope-sentry.top/",target:"_blank"},q={style:{position:"relative",top:"12px"}},ee={key:0},ae={class:"flex flex-col gap-2"};function le(e){return"function"==typeof e||"[object Object]"===Object.prototype.toString.call(e)&&!w(e)}const te=a({__name:"plugin",setup(a){const w=W({icon:"iconoir:search"}),{t:te}=C(),oe=l(""),ie=()=>{ve()},ne=t([{field:"index",label:te("tableDemo.index"),type:"index",minWidth:"15"},{field:"selection",type:"selection",minWidth:55},{field:"name",label:te("plugin.name"),formatter:(e,a,l)=>o("a",{href:`https://plugin.scope-sentry.top/plugin/${e.hash}`,style:"color: #409EFF; text-decoration: none;",target:"_blank"},[l])},{field:"module",label:te("plugin.module"),formatter:(e,a,l)=>{const t=se[l]||"#FFFFFF";return o(U,{style:{backgroundColor:t,color:"#000"}},le(l)?l:{default:()=>[l]})}},{field:"isSystem",label:te("plugin.isSystem"),formatter:(e,a,l)=>o(U,{type:l?"success":"warning"},{default:()=>[l?"true":"false"]})},{field:"version",label:te("plugin.version"),minWidth:100},{field:"parameter",label:te("plugin.parameter"),formatter:(e,a,l)=>o(B,{content:e.help,placement:"top",effect:"light"},{default:()=>[o("span",{style:"cursor: pointer;"},[l])]})},{field:"introduction",label:te("plugin.introduction"),minWidth:200},{field:"action",label:te("tableDemo.action"),minWidth:"300",fixed:"right",formatter:(e,a,l)=>{let t,c,d;const m=i(H,{onCommand:a=>{switch(a){case"reinstall":$("all",e.hash,e.module);break;case"recheck":J("all",e.hash,e.module);break;case"uninstall":O("all",e.hash,e.module)}}},{default:()=>i(n,{style:{outline:"none",boxShadow:"none"}},(()=>[te("common.operation"),i(s,{},(()=>i(r)))])),dropdown:()=>i(P,null,(()=>[i(z,{command:"reinstall"},(()=>te("plugin.reInstall"))),i(z,{command:"recheck"},(()=>te("plugin.reCheck"))),i(z,{command:"uninstall"},(()=>te("plugin.uninstall")))]))});return o(p,null,[m,o(u,{type:"warning",style:{marginLeft:"10px"},onClick:()=>Le(e)},le(t=te("common.log"))?t:{default:()=>[t]}),o(u,{type:"success",onClick:()=>Ae(e.id)},le(c=te("common.edit"))?c:{default:()=>[c]}),o(u,{type:"danger",onClick:()=>Fe(e.hash,e.module),disabled:e.isSystem},le(d=te("common.delete"))?d:{default:()=>[d]})])}}]),se={TargetHandler:"#2243dda6",SubdomainScan:"#FF9B85",SubdomainSecurity:"#FFFFBA",PortScanPreparation:"#BAFFB3",PortScan:"#BAE1FF",AssetMapping:"#e3ffba",URLScan:"#D1BAFF",WebCrawler:"#FFABAB",DirScan:"#3ccde6",VulnerabilityScan:"#FF677D",AssetHandle:"#B2E1FF",PortFingerprint:"#ffb5e4",URLSecurity:"#FFE4BA",PassiveScan:"#A2DFF7"},{tableRegister:re,tableState:ue,tableMethods:pe}=R({fetchDataApi:async()=>{const{currentPage:e,pageSize:a}=ue,l=await G(oe.value,e.value,a.value);return{list:l.data.list,total:l.data.total}},immediate:!0}),{loading:ce,dataList:de,total:me,currentPage:ge,pageSize:fe}=ue;fe.value=20;const{getList:ve,getElTableExpose:ye}=pe;function he(){return{background:"var(--el-fill-color-light)"}}const be=l(!1);let _e=te("plugin.new");const xe=()=>{be.value=!1},je=async()=>{c({title:"Delete",draggable:!0}).then((async()=>{await we()}))},Fe=async(e,a)=>{c({title:"Delete",draggable:!0}).then((async()=>{await ke(e,a)}))},Se=l(!1),ke=async(e,a)=>{Se.value=!0;try{await M([{hash:e,module:a}]);Se.value=!1,ve()}catch(l){Se.value=!1,ve()}},we=async()=>{const e=await ye(),a=((null==e?void 0:e.getSelectionRows())||[]).map((e=>({hash:e.hash,module:e.module})));Se.value=!0;try{await M(a);Se.value=!1,ve()}catch(l){Se.value=!1,ve()}},Ce=async()=>{Ee.value="",be.value=!0},Ee=l(""),Ae=async e=>{Ee.value=e,_e=te("common.edit"),be.value=!0};d((()=>{De(),window.addEventListener("resize",De)}));const Ve=l(0),De=()=>{const e=window.innerHeight||document.documentElement.clientHeight;Ve.value=.7*e},Ue=l(!1),Be=()=>{Ue.value=!1},Ie=l(""),Pe=l(),ze=l(""),He=l(""),Le=async e=>{ze.value=e.module,He.value=e.hash;const a=await N(e.module,e.hash);Ie.value=a.logs,Ue.value=!0},Re=async()=>{await K(ze.value,He.value),Ie.value=""},We={Authorization:`${m().getToken}`},Te=l(),Me=e=>{Te.value.clearFiles();const a=e[0];Te.value.handleStart(a)},Ne=e=>{var a;200===e.code?E.success("Upload succes"):E.error(e.message),505==e.code&&localStorage.removeItem("plugin_key"),ve(),null==(a=Te.value)||a.clearFiles()},Oe=(e,a)=>{a.length>0&&Te.value.submit()},Je=l(!1),$e=l(""),Ge=async()=>{if($e.value){200==(await Q($e.value)).code&&(localStorage.setItem("plugin_key",$e.value),Je.value=!1)}};return(()=>{const e=localStorage.getItem("plugin_key");e||(Je.value=!0),$e.value=e})(),(a,l)=>{const t=g("Icon");return f(),v(p,null,[o(h(e),null,{default:y((()=>[o(h(A),null,{default:y((()=>[o(h(V),{span:1},{default:y((()=>[o(h(D),{class:"mx-1",style:{position:"relative",top:"8px"}},{default:y((()=>[b(_(h(te)("plugin.name"))+":",1)])),_:1})])),_:1}),o(h(V),{span:5},{default:y((()=>[o(h(x),{modelValue:oe.value,"onUpdate:modelValue":l[0]||(l[0]=e=>oe.value=e),placeholder:h(te)("common.inputText"),style:{height:"38px"}},null,8,["modelValue","placeholder"])])),_:1}),o(h(V),{span:5,style:{position:"relative",left:"16px"}},{default:y((()=>[o(h(n),{type:"primary",icon:h(w),style:{height:"100%"},onClick:ie},{default:y((()=>[b("Search")])),_:1},8,["icon"])])),_:1})])),_:1}),o(h(A),{gutter:16,class:"mt-4"},{default:y((()=>[o(h(V),{xs:24,sm:24,md:24,lg:24,xl:24},{default:y((()=>[j("div",Y,[o(h(u),{type:"primary",onClick:Ce},{default:y((()=>[b(_(h(te)("plugin.new")),1)])),_:1}),o(h(u),{type:"danger",loading:Se.value,onClick:je},{default:y((()=>[b(_(h(te)("plugin.delete")),1)])),_:1},8,["loading"]),j("a",Z,[o(h(u),{type:"info"},{default:y((()=>[b(_(h(te)("plugin.market")),1)])),_:1})]),o(h(I),{ref_key:"upload",ref:Te,class:"flex items-center",action:"/api/plugin/import?key="+$e.value,headers:We,"on-success":Ne,limit:1,"on-exceed":Me,"auto-upload":!1,onChange:Oe},{trigger:y((()=>[o(h(u),null,{icon:y((()=>[o(t,{icon:"iconoir:upload"})])),default:y((()=>[b(" "+_(h(te)("plugin.import")),1)])),_:1})])),_:1},8,["action"])])])),_:1})])),_:1}),j("div",q,[o(h(L),{pageSize:h(fe),"onUpdate:pageSize":l[1]||(l[1]=e=>F(fe)?fe.value=e:null),currentPage:h(ge),"onUpdate:currentPage":l[2]||(l[2]=e=>F(ge)?ge.value=e:null),columns:ne,data:h(de),stripe:"",border:!0,loading:h(ce),resizable:!0,pagination:{total:h(me),pageSizes:[20,30,50,100,200,500,1e3]},onRegister:h(re),headerCellStyle:he,style:{fontFamily:"-apple-system,BlinkMacSystemFont,Segoe UI,Roboto,Helvetica Neue,Arial,Noto Sans,sans-serif,Apple Color Emoji,Segoe UI Emoji,Segoe UI Symbol,Noto Color Emoji"}},null,8,["pageSize","currentPage","columns","data","loading","pagination","onRegister"])])])),_:1}),o(h(T),{modelValue:be.value,"onUpdate:modelValue":l[3]||(l[3]=e=>be.value=e),title:h(_e),center:"",style:{"border-radius":"15px","box-shadow":"5px 5px 10px rgba(0, 0, 0, 0.3)"}},{default:y((()=>[o(X,{closeDialog:xe,getList:h(ve),id:Ee.value},null,8,["getList","id"])])),_:1},8,["modelValue","title"]),o(h(T),{modelValue:Ue.value,"onUpdate:modelValue":l[4]||(l[4]=e=>Ue.value=e),title:h(te)("node.log"),center:"",style:{"border-radius":"15px","box-shadow":"5px 5px 10px rgba(0, 0, 0, 0.3)"},maxHeight:Ve.value},{footer:y((()=>[o(h(u),{onClick:Re,type:"danger"},{default:y((()=>[b(_(h(te)("common.cleanLog")),1)])),_:1}),o(h(u),{onClick:Be},{default:y((()=>[b(_(h(te)("common.off")),1)])),_:1})])),default:y((()=>[o(h(S),{ref_key:"scrollbarRef",ref:Pe},{default:y((()=>[Ie.value?(f(),v("pre",ee,_(Ie.value),1)):k("",!0)])),_:1},512)])),_:1},8,["modelValue","title","maxHeight"]),o(h(T),{modelValue:Je.value,"onUpdate:modelValue":l[6]||(l[6]=e=>Je.value=e),title:h(te)("plugin.key"),center:"",width:"30%",style:{"max-width":"400px",height:"200px"}},{default:y((()=>[j("div",ae,[o(h(B),{class:"item",effect:"dark",content:h(te)("plugin.keyMsg"),placement:"top"},{default:y((()=>[o(h(x),{modelValue:$e.value,"onUpdate:modelValue":l[5]||(l[5]=e=>$e.value=e)},null,8,["modelValue"])])),_:1},8,["content"]),o(h(u),{onClick:Ge,type:"primary",class:"w-full"},{default:y((()=>[b("确定")])),_:1})])])),_:1},8,["modelValue","title"])],64)}}});export{te as default};
