import{_ as e}from"./ContentWrap.vue_vue_type_script_setup_true_lang-DYamwFJW.js";import{d as l,s as t,r as a,C as o,D as s,o as i,c as r,e as u,w as p,a as n,y as m,t as d,L as c,M as j,f as v,i as _,j as f,G as g,F as y,S as x,l as b,J as w,v as h}from"./index-vMt5tOuI.js";import{E as k,a as V}from"./el-tab-pane-zqlH3GjE.js";import{a as A,E as D}from"./el-col-CK66atlO.js";import{E}from"./el-switch-GToX5HcV.js";import{E as P}from"./el-text-JTnFT2W3.js";import{E as L,a as R}from"./el-form-CrZErb9v.js";import{E as S,a as U}from"./el-dropdown-item-C2usTBce.js";import"./el-popper-CJZHtU_A.js";import{_ as C}from"./ProjectList.vue_vue_type_style_index_0_lang-CIJmTnS2.js";import{_ as O}from"./AddProject.vue_vue_type_script_setup_true_lang-2JPxvtWN.js";import{_ as T}from"./Dialog.vue_vue_type_style_index_0_lang-BoLcmuNj.js";import{a as I,d as M}from"./index-DY_7WVuJ.js";import{u as $}from"./useIcon-D1gqy2ud.js";import"./el-card-CiG9T01-.js";import"./el-tooltip-l0sNRNKZ.js";import"./useInput-CPyPMgc8.js";import"./index-DMmJY3Ph.js";import"./isEqual-D2Iw95Gj.js";import"./debounce-heHWUU1E.js";import"./castArray-DsRHGUHq.js";import"./el-select-oq3jOrDj.js";import"./el-tag-C3mt1pWb.js";import"./refs-B-F76yq7.js";import"./Table.vue_vue_type_script_lang-kKqRQLK4.js";import"./el-table-column-Ty7HKRWd.js";import"./el-checkbox-DZAdOJwr.js";import"./isArrayLikeObject-4WqnSTxS.js";import"./raf-RUFdlNuH.js";import"./el-image-viewer-CkZtWhln.js";import"./tsxHelper-C3yG_Ynk.js";import"./index-BjDunm1X.js";/* empty css                          */import"./el-divider-RjUq80nQ.js";import"./el-radio-group-PFBUXVzP.js";import"./el-select-v2-iqUJrQ5P.js";import"./el-input-number-CYT48m1q.js";import"./index-Cei3qRzy.js";import"./index-znw49gEj.js";import"./index-CuGi13yK.js";import"./index-CcdpSLmp.js";const q={class:"mb-10px"},z=l({__name:"Project",setup(l){const{t:z}=b();let N=t({}),W=a([]),Y=t({});const B=a(!1),F=async(e,l)=>{0===e?(e=Q.value,l=Z.value):(Q.value=e,Z.value=l);try{const t=await I(X.value,e,l);Object.assign(N,t.data.result),W.value=Object.keys(t.data.tag),Object.assign(Y,t.data.tag);const a=W.value.indexOf("All");-1!==a&&W.value.splice(a,1)}catch(t){}},G=a(!1),H=async()=>{G.value=!0},J=()=>{G.value=!1},X=a(""),K=$({icon:"iconoir:search"}),Q=a(1),Z=a(50),ee=a(!1),le=async()=>{ee.value=!0,B.value=!0,await F(Q.value,Z.value),ee.value=!1,B.value=!1};le();const te=a(!1),ae=$({icon:"openmoji:delete"}),oe=async()=>{const e=a(!1);w({title:"Delete",draggable:!0,message:()=>h("div",{style:{display:"flex",alignItems:"center"}},[h("p",{style:{margin:"0 10px 0 0"}},z("task.delAsset")),h(E,{modelValue:e.value,"onUpdate:modelValue":l=>{e.value=l}})])}).then((async()=>{await M(ie.value,e.value),F(Q.value,Z.value)}))},se=$({icon:"ri:arrow-drop-down-line"}),ie=a([]);return(l,t)=>{const a=o("ElIcon"),b=o("ElDropdownMenu"),w=s("loading");return i(),r(y,null,[u(n(e),null,{default:p((()=>[u(n(A),{style:{"margin-bottom":"20px"},gutter:20},{default:p((()=>[u(n(D),{span:.5},{default:p((()=>[u(n(P),{class:"mx-1",style:{position:"relative",top:"8px"}},{default:p((()=>[m(d(n(z)("form.input"))+":",1)])),_:1})])),_:1}),u(n(D),{span:5},{default:p((()=>[u(n(c),{modelValue:X.value,"onUpdate:modelValue":t[0]||(t[0]=e=>X.value=e),placeholder:n(z)("common.inputText"),style:{height:"38px"}},null,8,["modelValue","placeholder"])])),_:1}),u(n(D),{span:5,style:{position:"relative",left:"16px"}},{default:p((()=>[u(n(j),{loading:ee.value,type:"primary",icon:n(K),size:"large",style:{height:"100%"},onClick:le},null,8,["loading","icon"])])),_:1})])),_:1}),u(n(A),{style:{"margin-bottom":"0%"}},{default:p((()=>[u(n(D),{span:2},{default:p((()=>[v("div",q,[u(n(j),{type:"primary",onClick:H},{default:p((()=>[m(d(n(z)("project.addProject")),1)])),_:1})])])),_:1}),u(n(D),{span:2},{default:p((()=>[u(n(L),null,{default:p((()=>[u(n(R),{label:n(z)("common.multipleSelection")},{default:p((()=>[u(n(E),{modelValue:te.value,"onUpdate:modelValue":t[1]||(t[1]=e=>te.value=e),class:"mb-2","inline-prompt":"","active-text":"Yes","inactive-text":"No"},null,8,["modelValue"])])),_:1},8,["label"])])),_:1})])),_:1}),te.value?(i(),_(n(D),{key:0,span:1},{default:p((()=>[u(n(S),{trigger:"click"},{dropdown:p((()=>[u(b,null,{default:p((()=>[u(n(U),{icon:n(ae),onClick:oe},{default:p((()=>[m(d(n(z)("common.delete")),1)])),_:1},8,["icon"])])),_:1})])),default:p((()=>[u(n(j),{plain:"",class:"custom-button align-bottom"},{default:p((()=>[m(d(n(z)("common.operation"))+" ",1),u(a,{class:"el-icon--right"},{default:p((()=>[u(n(se))])),_:1})])),_:1})])),_:1})])),_:1})):f("",!0)])),_:1}),g((i(),_(n(V),{class:"demo-tabs"},{default:p((()=>[u(n(k),{label:`All (${n(Y).All})`},{default:p((()=>[u(C,{tableDataList:n(N).All,getProjectTag:F,total:n(Y).All,multipleSelection:te.value,selectedRows:ie.value,"onUpdate:selectedRows":t[2]||(t[2]=e=>ie.value=e)},null,8,["tableDataList","total","multipleSelection","selectedRows"])])),_:1},8,["label"]),(i(!0),r(y,null,x(n(W),(e=>(i(),_(n(k),{label:`${e} (${n(Y)[e]})`,key:e},{default:p((()=>[u(C,{tableDataList:n(N)[e],getProjectTag:F,total:n(Y)[e],multipleSelection:te.value,selectedRows:ie.value,"onUpdate:selectedRows":t[3]||(t[3]=e=>ie.value=e)},null,8,["tableDataList","total","multipleSelection","selectedRows"])])),_:2},1032,["label"])))),128))])),_:1})),[[w,B.value]])])),_:1}),u(n(T),{modelValue:G.value,"onUpdate:modelValue":t[4]||(t[4]=e=>G.value=e),title:n(z)("project.addProject"),center:"",style:{"border-radius":"15px","box-shadow":"5px 5px 10px rgba(0, 0, 0, 0.3)"}},{default:p((()=>[u(O,{closeDialog:J,projectid:"",getProjectData:F,schedule:!1})])),_:1},8,["modelValue","title"])],64)}}});export{z as default};