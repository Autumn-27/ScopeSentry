const assert=require('node:assert/strict'),fs=require('node:fs'),path=require('node:path'),vm=require('node:vm'),{createRequire}=require('node:module');
const root=path.resolve(__dirname,'..');
const req=createRequire(path.join(root,'package.json')),ts=req('typescript'),{createI18n}=req('vue-i18n');
const messages=Object.fromEntries(['en','zh-CN','ko-KR'].map(lang=>[lang,vm.runInNewContext('('+fs.readFileSync(path.join(root,'src/locales',lang+'.ts'),'utf8').replace(/^export default/,'')+')')]));
const i18n=createI18n({legacy:false,locale:'ko-KR',fallbackLocale:'en',messages,missingWarn:false,fallbackWarn:false});
function load(file,mocks){const module={exports:{}};vm.runInNewContext(ts.transpileModule(fs.readFileSync(path.join(root,file),'utf8'),{compilerOptions:{module:ts.ModuleKind.CommonJS,target:ts.ScriptTarget.ES2020}}).outputText,{module,exports:module.exports,require:n=>n in mocks?mocks[n]:req(n)});return module.exports;}
const helper=load('src/axios/translateMessage.ts',{'@/plugins/vueI18n':{i18n}});
const calls=[];let logouts=0;
const {defaultResponseInterceptors:intercept}=load('src/axios/config.ts',{
 './translateMessage':helper,
 'element-plus':{ElMessage:{success:m=>calls.push(['success',m]),error:m=>calls.push(['error',m])}},
 '@/constants':{SUCCESS_CODE:200},'@/store/modules/user':{useUserStoreWithOut:()=>({logout:()=>logouts++})}
});
for(const lang of Object.keys(messages)){
 i18n.global.locale.value=lang;
 for(const key of ['api.success','api.system.version.success']){
  assert.equal(helper.translateApiMessage(key),i18n.global.t(key));
 }
 for(const raw of ['api.unknown.success','Normal server sentence','<script>alert(1)</script>','success'])assert.equal(helper.translateApiMessage(raw),raw);
 const data={code:200,message:'api.success',data:{unchanged:true}},response={headers:{},data};
 assert.equal(intercept(response),data);assert.equal(data.message,'api.success');assert.deepEqual(calls.pop(),['success',i18n.global.t('api.success')]);
 const failed={headers:{},data:{code:500,message:'api.success'}};
 assert.equal(intercept(failed),failed.data);assert.deepEqual(calls.pop(),['error',i18n.global.t('api.success')]);
}
i18n.global.setLocaleMessage('ko-KR',{...messages['ko-KR'],api:{}});i18n.global.locale.value='ko-KR';
assert.equal(helper.translateApiMessage('api.success'),messages.en.api.success);
for(const raw of [undefined,null,42,{}])assert.equal(helper.translateApiMessage(raw),'');
const blob={headers:{},config:{responseType:'blob'},data:{}};assert.equal(intercept(blob),blob);
const unauthorized={headers:{},data:{code:401,message:'Unknown error'}};assert.equal(intercept(unauthorized),unauthorized.data);assert.equal(logouts,1);
console.log('PASS: registered API keys in three locales; English fallback; raw/unknown/HTML text preserved; status channel, response identity, blob bypass and 401 behavior unchanged.');
