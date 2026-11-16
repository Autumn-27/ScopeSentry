# ScopeSentry Kubernetes 部署说明

本目录是把 ScopeSentry 拆成 Kubernetes 资源后的部署清单，对应 docker compose 版本（`../single-host-deployment.yml`）的四容器结构：

| 组件 | 资源 | 镜像 | 说明 |
| --- | --- | --- | --- |
| 命名空间 | `Namespace/scope-sentry` | – | 所有资源都落在这里 |
| MongoDB | `StatefulSet/mongodb` + `Service/mongodb` | `mongo:7.0.28` | 存业务数据，**需要持久卷** |
| Redis | `Deployment/redis` + `Service/redis` | `redis:7.0.11` | 任务队列/缓存，**重启即丢** |
| 服务端 | `Deployment/webserver` + `Service/scopesentry`(NodePort) | `autumn27/scopesentry:1.9.4` | Web/API，容器端口 8082 |
| 扫描端 | `Deployment/scan` | `autumn27/scopesentry-scan:1.9.1` | 无监听端口，从 Redis 取任务 |

目录结构：

```
k8s/
├── all.yml                        # 上述 10 个资源的合并版（与分文件内容完全一致）
├── namespace/namespace.yml
├── secrets/mongo-secret.yml
├── secrets/redis-secret.yml
├── pods/mongodb-statefulset.yml
├── pods/redis-deployment.yml
├── pods/scope-sentry-deployment.yml
├── pods/scan-deployment.yml
├── services/mongodb-service.yml
├── services/redis-service.yml
└── services/scope-sentry-service.yml
```

`all.yml` 是分文件的拼接版（已逐文档比对，10/10 一致）。**改清单时两处都要改**，否则会漂移；CI 可加校验见文末。

---

## 一、校验方法与当前结果

在仓库根目录执行。以下输出为 2026-11-16 在本分支（`docs/k8s-validate`）实测。

### 1. YAML 语法解析

```bash
python3 -c "import yaml,glob,sys; [list(yaml.safe_load_all(open(f))) for f in glob.glob('k8s/**/*.yml', recursive=True)]; print('yaml ok')"
```

```
yaml ok
```

逐文件文档数（确认没有空文档、没有被 `---` 切碎）：

```
k8s/all.yml                            docs=10 kinds=['Namespace','Secret','Secret','StatefulSet','Deployment','Service','Service','Deployment','Service','Deployment']
k8s/namespace/namespace.yml            docs= 1 kinds=['Namespace']
k8s/pods/mongodb-statefulset.yml       docs= 1 kinds=['StatefulSet']
k8s/pods/redis-deployment.yml          docs= 1 kinds=['Deployment']
k8s/pods/scan-deployment.yml           docs= 1 kinds=['Deployment']
k8s/pods/scope-sentry-deployment.yml   docs= 1 kinds=['Deployment']
k8s/secrets/mongo-secret.yml           docs= 1 kinds=['Secret']
k8s/secrets/redis-secret.yml           docs= 1 kinds=['Secret']
k8s/services/mongodb-service.yml       docs= 1 kinds=['Service']
k8s/services/redis-service.yml         docs= 1 kinds=['Service']
k8s/services/scope-sentry-service.yml  docs= 1 kinds=['Service']
```

### 2. kubectl 校验（无集群时的行为，重要）

```bash
kubectl apply --dry-run=client -f k8s/all.yml
```

当前机器上 **没有可用集群、也没有 kubeconfig**（`~/.kube` 为空，`kubectl config get-contexts` 无输出），因此：

```
error: error validating "k8s/all.yml": error validating data: failed to download openapi: Get "http://localhost:8080/openapi/v2?timeout=32s": dial tcp [::1]:8080: connect: connection refused
```

注意：`--dry-run=client` **并不意味着不需要 API server**。kubectl 的字段校验和资源类型识别都要先拉 `/openapi/v2` 与 discovery 信息；加 `--validate=false` 也只是跳过 schema 校验，仍会卡在 `unable to recognize "k8s/all.yml"`（`kubectl create --dry-run=client` 同样如此）。所以「无集群跑 kubectl dry-run」这条路走不通，必须用离线 schema 校验器替代。

真正能给出结论的做法：在**有集群的环境**执行

```bash
kubectl apply --dry-run=client -f k8s/all.yml     # 1.27+ 默认 server-side field validation=Strict
kubectl apply --dry-run=server -f k8s/all.yml     # 需要真实集群，可额外发现 CRD/准入校验问题
```

其中 `namespace 不存在` 的报错要单独区分：本清单在 `all.yml` 里第一个文档就是 `Namespace`，按文件顺序 apply 不会出现该问题；但若只 apply 某个分文件（例如 `kubectl apply -f k8s/pods/`），会报 `Error from server (NotFound): namespaces "scope-sentry" not found` —— 这不是清单错误，是执行顺序问题。

### 3. kubeconform 离线 schema 校验（本次实际采用）

```bash
# 安装（macOS arm64 示例）
curl -sSL https://github.com/yannh/kubeconform/releases/download/v0.6.7/kubeconform-darwin-arm64.tar.gz | tar -xz -C /tmp

/tmp/kubeconform -strict -summary k8s/all.yml
```

```
k8s/all.yml - Namespace scope-sentry is valid
k8s/all.yml - Secret mongo-secret is valid
k8s/all.yml - Secret redis-secret is valid
k8s/all.yml - Deployment redis is valid
k8s/all.yml - Deployment webserver is valid
k8s/all.yml - Service mongodb is valid
k8s/all.yml - Deployment scan is valid
k8s/all.yml - Service redis is valid
k8s/all.yml - Service scopesentry is valid
k8s/all.yml - StatefulSet mongodb is invalid: problem validating schema. Check JSON formatting: jsonschema: '' does not validate with .../master-standalone-strict/statefulset-apps-v1.json#/additionalProperties: additionalProperties 'volumeClaimTemplates' not allowed
Summary: 10 resources found in 1 file - Valid: 9, Invalid: 1, Errors: 0, Skipped: 0
```

全部文件一起跑（结果一致，说明 `all.yml` 与分文件是同一份内容）：

```bash
/tmp/kubeconform -strict -summary $(find k8s -name '*.yml' | sort)
```

```
k8s/all.yml - StatefulSet mongodb is invalid: ... additionalProperties 'volumeClaimTemplates' not allowed
k8s/pods/mongodb-statefulset.yml - StatefulSet mongodb is invalid: ... additionalProperties 'volumeClaimTemplates' not allowed
Summary: 20 resources found in 11 files - Valid: 18, Invalid: 2, Errors: 0, Skipped: 0
```

固定 Kubernetes 版本（1.31.0）复跑结论相同。**唯一的 schema 错误就是下面第 1 号阻断问题**，其余 18 个资源全部合法。

### 4. 引用一致性静态核查

本次用脚本核对（可直接复现）：

```bash
python3 - <<'PY'
import yaml, glob
docs=[]
for f in sorted(glob.glob('k8s/**/*.yml', recursive=True)):
    if f=='k8s/all.yml': continue
    docs += [d for d in yaml.safe_load_all(open(f)) if d]
secrets={}
for d in docs:
    if d['kind']=='Secret':
        secrets[d['metadata']['name']] = set((d.get('stringData') or {})) | set((d.get('data') or {}))
print("secrets:", secrets)
bad=0
for d in docs:
    if d['kind'] not in ('Deployment','StatefulSet'): continue
    for c in d['spec']['template']['spec']['containers']:
        for e in c.get('env') or []:
            r=(e.get('valueFrom') or {}).get('secretKeyRef')
            if r and not (r['name'] in secrets and r['key'] in secrets[r['name']]):
                bad+=1; print("BROKEN", d['metadata']['name'], e['name'], r)
print("broken secretKeyRef:", bad)
PY
```

```
secrets: {'mongo-secret': {'MONGO_INITDB_ROOT_PASSWORD','MONGO_INITDB_ROOT_USERNAME'},
          'redis-secret': {'REDIS_PASSWORD'}}
broken secretKeyRef: 0
```

| 检查项 | 结果 |
| --- | --- |
| `secretKeyRef` → Secret/key 是否都定义 | ✅ 9 处引用全部命中，无悬空 |
| `Service.selector` → Pod labels | ✅ `mongodb→app=mongodb`、`redis→app=redis`、`scopesentry→app=webserver` 全部命中唯一 workload |
| `StatefulSet.serviceName` → Service 名 | ✅ `mongodb` 存在（但该 Service 不是 headless，见 9 号） |
| 镜像 tag ↔ `version.json` | ✅ server `1.9.4` / scan `1.9.1`，与 `version.json` 完全一致；4 个镜像 tag 在 Docker Hub 均存在 |
| 探针 / resources / PVC | ⚠️ 4 个 workload 全无探针、全无 resources；Mongo 的 PVC 定义失效（1 号）；其余 3 个无持久卷 |
| NodePort 冲突 | ✅ `30082` 在默认 30000–32767 内，不与任何知名服务端口冲突（详见 20 号） |

---

## 二、部署顺序

顺序不能乱：`namespace → secrets → mongo/redis → server/scan`。

```bash
# 0) 先确认集群有默认 StorageClass，否则 Mongo 的 PVC 会一直 Pending
kubectl get storageclass

# 1) 命名空间
kubectl apply -f k8s/namespace/namespace.yml

# 2) 密钥（建议按第三节替换成真实口令后再 apply）
kubectl apply -f k8s/secrets/mongo-secret.yml -f k8s/secrets/redis-secret.yml

# 3) 中间件：先 Service 再工作负载（StatefulSet 的 serviceName 依赖 Service 存在）
kubectl apply -f k8s/services/mongodb-service.yml -f k8s/services/redis-service.yml
kubectl apply -f k8s/pods/mongodb-statefulset.yml -f k8s/pods/redis-deployment.yml

# 4) 等中间件真正可用再做下一步（没有探针，kubectl wait 只能等 Pod Ready，不够）
kubectl -n scope-sentry rollout status statefulset/mongodb --timeout=180s
kubectl -n scope-sentry rollout status deployment/redis   --timeout=120s
kubectl -n scope-sentry exec sts/mongodb -- mongosh --quiet --eval 'db.adminCommand({ping:1})'

# 5) 应用层
kubectl apply -f k8s/pods/scope-sentry-deployment.yml -f k8s/services/scope-sentry-service.yml
kubectl apply -f k8s/pods/scan-deployment.yml

# 6) 一条命令全量（all.yml 里 Namespace 是第一个文档，顺序安全）
kubectl apply -f k8s/all.yml
```

⚠️ `kubectl apply -f k8s/`（不带 `-R`）会报 `error: no objects passed to apply`，因为 `k8s/` 根目录下没有 `*.yml`，只有子目录。要么用 `-f k8s/all.yml`，要么用 `-Rf k8s/`（后者会把 `all.yml` 和分文件各 apply 一遍，幂等但啰嗦）。

验证：

```bash
kubectl -n scope-sentry get pod,svc,pvc -o wide
kubectl -n scope-sentry rollout status deployment/webserver --timeout=180s
kubectl -n scope-sentry logs deploy/webserver --tail=50 | grep -Ei 'mongo|redis|version|Starting server'
```

启动成功的标志（服务端日志）：

```
MongoDB connected successfully
Redis connected successfully
version: 1.9.4
Starting server  {"port": 8082}
```

---

## 三、Secret 填写与替换

两个 Secret 都用 `stringData`（明文写在清单里，apply 时由 API server 转成 base64）。键名**必须**与下列一致，否则 Pod 起不来（`CreateContainerConfigError`）：

| Secret | 键 | 被谁引用 |
| --- | --- | --- |
| `mongo-secret` | `MONGO_INITDB_ROOT_USERNAME` | mongo 初始化、webserver/scan 的 `MONGODB_USER` |
| `mongo-secret` | `MONGO_INITDB_ROOT_PASSWORD` | mongo 初始化、webserver/scan 的 `MONGODB_PASSWORD` |
| `redis-secret` | `REDIS_PASSWORD` | redis `--requirepass`、webserver/scan 的 `REDIS_PASSWORD` |

### 方式 A：改文件（团队仓库场景，注意别把真口令提交上去）

```yaml
stringData:
  MONGO_INITDB_ROOT_USERNAME: admin
  MONGO_INITDB_ROOT_PASSWORD: <换成强口令>
```

### 方式 B：命令行生成（推荐，不落盘）

```bash
kubectl -n scope-sentry create secret generic mongo-secret \
  --from-literal=MONGO_INITDB_ROOT_USERNAME=admin \
  --from-literal=MONGO_INITDB_ROOT_PASSWORD="$(openssl rand -base64 24)" \
  --dry-run=client -o yaml | kubectl apply -f -

kubectl -n scope-sentry create secret generic redis-secret \
  --from-literal=REDIS_PASSWORD="$(openssl rand -base64 24)" \
  --dry-run=client -o yaml | kubectl apply -f -
```

### 三个必须知道的约束

1. **Mongo 口令只在首次初始化时生效**。`MONGO_INITDB_ROOT_*` 只在 `/data/db` 为空时执行；PVC 里已有数据时改 Secret 不会改数据库里的口令，只会让 webserver/scan 用新口令去连旧库 → 认证失败。改口令要走 `db.changeUserPassword`，或删 PVC 重来（会丢数据）。
2. **改完 Secret 要重启消费者**，env 不会热更新：
   ```bash
   kubectl -n scope-sentry rollout restart deploy/webserver deploy/scan statefulset/mongodb deploy/redis
   ```
3. **`data` 与 `stringData` 不要混用同名键**，否则以 `data` 为准（base64 解不开就静默变空串）。当前清单统一用 `stringData`，保持一致即可。

读取校验：

```bash
kubectl -n scope-sentry get secret mongo-secret -o jsonpath='{.data.MONGO_INITDB_ROOT_USERNAME}' | base64 -d; echo
```

---

## 四、修改镜像 tag

当前 tag 与 `version.json`（仓库根目录）一致：

```json
{ "server": "1.9.4", "scan": "1.9.1" }
```

- 服务端：`k8s/pods/scope-sentry-deployment.yml` 与 `k8s/all.yml`（`Deployment/webserver`）里的 `image: autumn27/scopesentry:1.9.4`
- 扫描端：`k8s/pods/scan-deployment.yml` 与 `k8s/all.yml`（`Deployment/scan`）里的 `image: autumn27/scopesentry-scan:1.9.1`

升级步骤：

```bash
# 1) 改 tag（两处：分文件 + all.yml，务必一致）
sed -i '' 's|autumn27/scopesentry:1.9.4|autumn27/scopesentry:1.9.5|' k8s/pods/scope-sentry-deployment.yml k8s/all.yml

# 2) 先确认 tag 真的存在，避免 ImagePullBackOff
curl -s -o /dev/null -w '%{http_code}\n' https://hub.docker.com/v2/repositories/autumn27/scopesentry/tags/1.9.5   # 期望 200

# 3) apply + 观察
kubectl apply -f k8s/pods/scope-sentry-deployment.yml
kubectl -n scope-sentry rollout status deploy/webserver --timeout=180s
```

建议直接用 `kubectl set image` 做应急回滚/升级（不需要改文件，但会和仓库清单产生偏离，记得回写）：

```bash
kubectl -n scope-sentry set image deploy/webserver webserver=autumn27/scopesentry:1.9.5
kubectl -n scope-sentry rollout undo deploy/webserver
```

注意事项：

- 不要用 `:latest`。清单里全是精确 tag，`imagePullPolicy` 默认 `IfNotPresent`，行为可预期；换成 `latest` 会变成 `Always` 且无法回滚定位。
- **服务端与扫描端版本要匹配**。`version.json` 把两者分开标版本（1.9.4 / 1.9.1），升级服务端时确认扫描端兼容；扫描端源码在独立仓库 `Autumn-27/ScopeSentry-Scan`。
- 服务端启动时会跑数据库迁移（`internal/update/update.go`，按 Mongo 里 `config.version` 判断），跨多个版本升级时**先升服务端、等迁移跑完再升扫描端**，否则新旧节点可能按不同 schema 读写。
- `imagePullPolicy` 若要强制每次拉取（例如内网 registry 覆盖同名 tag），显式写 `imagePullPolicy: Always`。

---

## 五、如何暴露服务

现状：`Service/scopesentry` 是 `NodePort 30082 → 8082`，即 `http://<任一节点IP>:30082`。可用，但有两个前提：

1. 节点安全组/防火墙放通 30082（很多云默认只开 80/443，30000+ 需要手动加）。
2. 集群的 `--service-node-port-range` 包含 30082（kubeadm/k3s/RKE2 默认 30000–32767，均满足）。

`8082` 是服务端默认端口（`internal/config/config.go` 的 `server.port` 默认值），监听地址是 `0.0.0.0`（`cmd/main/main.go` 里 `Addr: ":8082"`），所以 NodePort/targetPort 都能通。

### 方案 1：保持 NodePort（最简单）

```bash
kubectl -n scope-sentry get svc scopesentry -o wide
# 拿到 30082 后直接访问
curl -I http://<NODE_IP>:30082/
```

换端口时只改 `nodePort`（30000–32767），`port`/`targetPort` 保持 8082：

```yaml
ports:
  - port: 8082
    targetPort: 8082
    nodePort: 30082
    protocol: TCP
```

### 方案 2：Ingress + TLS（生产建议）

NodePort 是明文 HTTP，登录口令和 token 裸奔。生产建议把 Service 改回 `ClusterIP`，用 Ingress 收口：

```yaml
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: scopesentry
  namespace: scope-sentry
  annotations:
    nginx.ingress.kubernetes.io/proxy-body-size: "200m"   # 插件 zip / 字典上传
    nginx.ingress.kubernetes.io/proxy-read-timeout: "300"
spec:
  ingressClassName: nginx
  tls:
    - hosts: [sentry.example.com]
      secretName: sentry-tls
  rules:
    - host: sentry.example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name: scopesentry
                port:
                  number: 8082
```

前端是 SPA（`NoRoute` 统一回 `index.html`），并且 `/api`、`/assets`、`/js`、`/plugins`、`/images` 都是同源相对路径，**必须把 `/` 整个域名交给它，不要加 path 前缀重写**，否则静态资源 404。

### 方案 3：临时本地访问

```bash
kubectl -n scope-sentry port-forward svc/scopesentry 8082:8082
# 浏览器打开 http://127.0.0.1:8082
```

---

## 六、环境变量对照表

服务端（本仓库）读取的变量在 `internal/config/config.go: loadFromEnv()`，扫描端在独立仓库。对照 compose 可以看清哪些变量在 k8s 里漏了：

| 变量 | 用途 | compose(server) | compose(scan) | k8s(webserver) | k8s(scan) |
| --- | --- | --- | --- | --- | --- |
| `TIMEZONE` | 服务端时区（本仓库唯一读取的时区变量） | ✅ | – | ✅ | ✅（但扫描镜像未必读） |
| `TimeZoneName` | 扫描端时区（compose 给 scan 用的键） | – | ✅ | – | ❌ **缺** |
| `NodeName` | 扫描节点名 | – | ✅ `node-test` | – | ❌ **缺** |
| `MONGODB_IP` | Mongo 地址 | ✅ | ✅ | ✅ `mongodb` | ✅ `mongodb` |
| `MONGODB_PORT` | Mongo 端口 | ✅ | ✅ | ✅ `27017` | ✅ `27017` |
| `MONGODB_DATABASE` | 库名 | ✅ | ✅ | ✅ `ScopeSentry` | ✅ `ScopeSentry` |
| `MONGODB_USER` | 用户名 | ✅ | ✅ | ✅ secretKeyRef | ✅ secretKeyRef |
| `MONGODB_PASSWORD` | 口令 | ✅ | ✅ | ✅ secretKeyRef | ✅ secretKeyRef |
| `REDIS_IP` | Redis 地址 | ✅ | ✅ | ✅ `redis` | ✅ `redis` |
| `REDIS_PORT` | Redis 端口 | ✅ | ✅ | ✅ `6379` | ✅ `6379` |
| `REDIS_PASSWORD` | Redis 口令 | ✅ | ✅ | ✅ secretKeyRef | ✅ secretKeyRef |

### NodeName 到底要不要？（按代码核查结论）

**本仓库（服务端）完全不读 `NodeName`** —— 全仓 `grep -rn NodeName` 只命中 `single-host-deployment.yml` 和 `single-host-deployment-fix.yml`，Go 代码零引用。

但扫描端是**独立仓库** `Autumn-27/ScopeSentry-Scan`（见 `README_CN.md:25`），其部署说明写得很明确（`README_CN.md:49`）：

> `NodeName` 为节点名称，每个节点名称请保持唯一性（**如果为空会随机生成**，可在 web 界面更改名称）

所以结论是：**k8s 的 scan Deployment 应该显式设 `NodeName`**。不设不会启动失败，但每次 Pod 重建都会以随机名字注册成一个「新节点」，后果是：

- Web「节点」页不断累积离线历史节点（服务端按 Redis `node:<name>` 键管理节点，见 `internal/repositories/node/node.go`）；
- 无法把扫描任务稳定绑定到某个节点；
- 节点超时判定（`NodeTimeout = 300s`）会把这些僵尸节点标成离线，干扰运维判断。

建议补上（单副本用固定值最省事；多副本用 downward API 取 Pod 名）：

```yaml
env:
  - name: NodeName
    valueFrom:
      fieldRef:
        fieldPath: metadata.name      # 每副本唯一；注意重启后 Pod 名会变
  # 或者固定值，适合 replicas: 1
  # - name: NodeName
  #   value: scan-1
  - name: TimeZoneName
    value: Asia/Shanghai              # 扫描镜像按 compose 的写法读这个键
```

`TimeZoneName` 同理：compose 给 scan 容器传的是 `TimeZoneName`（不是 `TIMEZONE`），k8s 里只给了 `TIMEZONE`。两个都设最稳妥，多传一个变量没有副作用。

---

## 七、常见坑

### 1. Mongo 的持久卷现在其实是「没有」

清单里写了 `volumeClaimTemplates` 20Gi，但**缩进错位**（见 1 号阻断问题），实际不会创建 PVC。所以当前状态等价于「单副本 Mongo 无 PVC」：`/data/db` 落在容器可写层，**Pod 重建 = 全库丢失**（项目、资产、任务、扫描结果全没）。修好缩进后才真正有 20Gi PVC，但仍有以下限制：

- `storageClassName` 没写 → 依赖集群默认 StorageClass；没有默认 SC 时 PVC 永远 `Pending`，Pod `Pending`。
- `accessModes: ReadWriteOnce` → 只能被单节点挂载，Pod 漂移到别的节点时依赖 CSI 的 attach/detach。
- 删 StatefulSet **不会**删 PVC（`persistentVolumeClaimRetentionPolicy` 默认为 `Retain`），重装时会复用到旧数据 —— 想彻底重来要手动 `kubectl delete pvc`。
- 单副本、无副本集：代码未使用 Mongo 事务（全仓无业务侧 `WithTransaction` 调用），所以单节点能用；但没有高可用，节点故障即停机。

### 2. 服务端每次重启都会踢掉所有登录用户

`internal/config/config.go` 在启动时 `viper.WriteConfigAs(exeDir/config.yaml)` **重写配置文件**，而 `jwt.secret` 的默认值是 `random.GenerateString(32)`。容器没有 PVC 时 `config.yaml` 随容器消失，下次启动就是新密钥 → **所有已签发的 token 立即失效**，用户全部被登出。

同一位置还会重新生成：

- `PLUGINKEY`（`handlePluginKey()`，写在 CWD，容器里是 `/opt/ScopeSentry/PLUGINKEY`）→ 插件密钥变化；
- `uploads/`（插件 zip、POC 文件，代码里用的是**相对路径** `uploads/xxx.zip`）、`images/screenshots`（扫描截图，`/images` 静态路由直接读盘）、`files/export`（导出文件）→ 全部丢失，历史截图变 404。

修法：给 `webserver` 挂一个 PVC 到 `/opt/ScopeSentry`（注意会覆盖镜像里的二进制，需要 `subPath` 或只挂子目录），或至少分别挂 `uploads`/`images`/`files` 三个子目录，并把 `config.yaml` 用 ConfigMap 固定下来（**关键是固定 `jwt.secret`**）。

### 3. 启动顺序没保障，服务端会 panic 重启

`internal/database/mongodb` 和 `internal/database/redis` 都是包级 `init()`：连接失败只打日志然后 `return`，`Client`/`DB` 保持 `nil`。而 `main()` 紧接着调用 `update.Update()`，里面直接 `mongodb.DB.Collection("config")` → **空指针 panic**，容器退出、重启。

k8s 里 4 个 Deployment 是并行创建的，Mongo 首次初始化还要十几秒，所以「服务端比 Mongo 先起来」几乎必然发生一次。表现就是 `webserver` / `scan` 先 CrashLoopBackOff 1~3 次再自己好（重启退避最长会到几十秒）。compose 版本用 `depends_on: condition: service_healthy` 回避了这个问题，k8s 没有对应机制。

缓解（任选）：

```yaml
# A. initContainer 等待依赖
initContainers:
  - name: wait-mongo
    image: busybox:1.36
    command: ['sh','-c','until nc -z mongodb 27017; do echo waiting mongo; sleep 2; done']
  - name: wait-redis
    image: busybox:1.36
    command: ['sh','-c','until nc -z redis 6379; do echo waiting redis; sleep 2; done']
```

```yaml
# B. 给中间件加 readinessProbe，让它自己先「不 Ready」
readinessProbe:                      # mongo
  exec: { command: ["mongosh","--quiet","--eval","db.adminCommand({ping:1})"] }
  initialDelaySeconds: 10
  periodSeconds: 5
readinessProbe:                      # redis（注意要带口令）
  exec: { command: ["sh","-c","redis-cli -a \"$REDIS_PASSWORD\" --no-auth-warning ping | grep -q PONG"] }
  initialDelaySeconds: 5
  periodSeconds: 5
```

initContainer 方案更彻底：它保证应用容器第一次启动时依赖已经可连。

### 4. 4 个 workload 都没有探针和资源限制

- 没有 `readinessProbe` → Pod 一启动就被标 Ready、立刻进 Service Endpoints。Mongo 首次初始化期间也会被当成就绪，流量打过去直接失败。
- 没有 `livenessProbe` → 进程假死（例如 Mongo 连接泄漏后 hang 住）不会被 k8s 重启。
- 没有 `resources.requests/limits` → 调度器按 0 资源算，能塞多少塞多少；`scan` 端跑 nmap/爬虫/截图时 CPU 和内存会尖峰，**没有 limit 就可能把整个节点拖垮**（OOM killer 杀别的 Pod，或磁盘被日志写满）。这是这套清单里风险最高的一条运维缺口。

建议基线（按实际规模调）：

```yaml
resources:
  requests: { cpu: "100m", memory: "256Mi" }
  limits:   { cpu: "2",    memory: "2Gi" }     # scan 需要更宽
```

### 5. `kubectl apply -f k8s/` 什么都不会创建

`k8s/` 根目录只有子目录，没有 `*.yml`。必须 `-f k8s/all.yml` 或 `-Rf k8s/`。

### 6. all.yml 与分文件是两份真相

`all.yml` 是拼接版，当前 10 个文档与分文件**逐文档一致**（已校验）。任何人只改一边就会漂移，而 CI 通常只校验 `all.yml`（因为一条命令能跑完），漂移会静默生效。建议加个 CI 校验（见附录）。

### 7. NodePort 30082 的三个隐性前提

- `nodePort` 在**集群范围内唯一**，若已被别的 Service 占用，apply 报 `provided port is already allocated`；
- 云厂商节点安全组默认常只放通 80/443，30082 要手动加；
- 明文 HTTP，登录口令与 token 不加密。生产用 Ingress + TLS（第五节方案 2）。

### 8. Redis 不是缓存，是队列

Redis 里放的是扫描任务队列和节点心跳（`node:<name>`）。当前 `Deployment/redis` 无持久卷、无 `strategy: Recreate`、无探针，滚动更新时旧 Pod 立刻被摘除、新 Pod 立刻 Ready → **重启期间队列全丢**，正在跑的任务消失。要么给它挂 PVC（`/data`），要么接受「重启即丢队列」并在升级时避开任务高峰。

---

## 八、已知问题与建议

### A. 阻断部署（必须先修，否则起不来）

| # | 问题 | 位置 | 影响 |
| --- | --- | --- | --- |
| 1 | `volumeClaimTemplates` 缩进错位，落在**文档根级**而不是 `spec` 下 | `k8s/pods/mongodb-statefulset.yml:45`、`k8s/all.yml:86` | kubeconform `-strict` 直接判 invalid（`additionalProperties 'volumeClaimTemplates' not allowed`）。kubectl ≥1.27 默认 server-side field validation=Strict，apply 会报 `strict decoding error: unknown field "volumeClaimTemplates"` 而整体失败；若用 `--validate=warn` 丢弃该字段，则 PVC 不会创建，而 Pod 模板里 `volumeMounts: mongodb-data` 找不到对应卷，StatefulSet 无法创建 Pod → **Mongo 永不就绪 → webserver/scan 连不上库 → 全站不可用** |
| 2 | 未指定 `storageClassName`（条件性阻断） | 同上 | 集群没有默认 StorageClass 时 PVC 一直 `Pending`，Pod `Pending`。部署前必须 `kubectl get sc` 确认，否则修好 1 号也起不来 |

**1 号修复补丁**（两处都要改，`volumeClaimTemplates` 与 `selector`/`template` 同级，缩进 2 空格）：

```diff
 spec:
   serviceName: mongodb
   replicas: 1
   selector:
     matchLabels:
       app: mongodb
   template:
     ...
-volumeClaimTemplates:
-  - metadata:
-      name: mongodb-data
-
-    spec:
-      accessModes:
-        - ReadWriteOnce
-
-      resources:
-        requests:
-          storage: 20Gi
+  volumeClaimTemplates:
+    - metadata:
+        name: mongodb-data
+        labels:
+          app: mongodb
+      spec:
+        accessModes:
+          - ReadWriteOnce
+        resources:
+          requests:
+            storage: 20Gi
+        # storageClassName: <集群实际 SC>   # 没有默认 SC 时必须显式指定
```

修完复验：`kubeconform -strict -summary k8s/all.yml` 应为 `Valid: 10, Invalid: 0`。

### B. 建议修复（能跑，但会咬人）

| # | 问题 | 位置 | 影响 / 建议 |
| --- | --- | --- | --- |
| 3 | 4 个 workload 全部没有 `readinessProbe`/`livenessProbe` | 全部 pods | 依赖未就绪即被判 Ready；进程假死不重启。补探针（第七节第 3、4 条给了可直接用的片段） |
| 4 | 4 个 workload 全部没有 `resources` | 全部 pods | 调度不可预期；scan 端 CPU/内存尖峰可拖垮节点。至少给 scan 加 limits |
| 5 | `webserver` 的 `/opt/ScopeSentry` 无持久卷 | `scope-sentry-deployment.yml` | **每次重启 JWT secret 重置 → 全员掉登录**；`PLUGINKEY`、`uploads/`、`images/screenshots`、`files/export` 全丢。挂 PVC 或固定 `config.yaml` |
| 6 | 无启动顺序保障，依赖未就绪时 `main()` 空指针 panic | `webserver`、`scan` | 必然先 CrashLoop 几次再自愈（重启退避可达数十秒）。加 initContainer 等 27017/6379 |
| 7 | `scan` 未设 `NodeName` | `scan-deployment.yml` | 每次重建随机节点名，Web 节点页累积僵尸节点，任务无法绑定固定节点。见第六节 |
| 8 | `scan` 只给了 `TIMEZONE`，compose 给扫描端的是 `TimeZoneName` | `scan-deployment.yml` | 扫描端时区可能没生效（影响日志/任务时间戳）。两个都设 |
| 9 | `Service/mongodb` 不是 headless，但 `StatefulSet.serviceName` 指向它 | `mongodb-service.yml` | `mongodb-0.mongodb` 不解析。当前应用只用 `mongodb:27017`（ClusterIP），能用；但做副本集/单 Pod 直连时会踩坑。建议另建 `clusterIP: None` 的 headless Service |
| 10 | `redis` 无持久卷、无 `strategy: Recreate` | `redis-deployment.yml` | 重启即丢任务队列；滚动更新期间队列清空 |
| 11 | `redis-server --requirepass $(REDIS_PASSWORD)` | `redis-deployment.yml` | 口令出现在进程 argv，容器内 `ps` 可见，也会进 `/proc/1/cmdline`。改用配置文件或 `REDIS_ARGS` |
| 12 | `all.yml` 与分文件双份维护 | `k8s/` | 当前一致，但无机制保证。加 CI 校验（附录） |
| 13 | 无 `securityContext`、无 `NetworkPolicy` | 全部 pods | 默认 root、可写根文件系统、Pod 间无隔离。至少 `runAsNonRoot` + `allowPrivilegeEscalation: false`（注意服务端会写 `/opt/ScopeSentry`，加 `readOnlyRootFilesystem` 前要先挂卷） |
| 14 | Secret 用 `stringData` 明文提交在仓库里，口令是 `mongodb_password`/`redis_password` | `secrets/*.yml` | 弱口令 + 随 git 分发。用第三节方式 B 生成，或接 SealedSecret/ExternalSecret |
| 15 | 暴露方式是明文 NodePort | `scope-sentry-service.yml` | 登录口令/token 不加密。生产上 Ingress + TLS（第五节方案 2） |

### C. 可接受（已核查，无需处理）

| # | 项 | 结论 |
| --- | --- | --- |
| 16 | 单副本 Mongo、非副本集 | 代码未使用 Mongo 事务（全仓无业务侧 `WithTransaction`），单节点可用。**前提是 1 号修好、PVC 真的创建出来** |
| 17 | 镜像 tag 与仓库版本 | `autumn27/scopesentry:1.9.4`、`autumn27/scopesentry-scan:1.9.1` 与 `version.json` 完全一致；4 个镜像 tag 在 Docker Hub 均存在（逐个查过 tags API，返回 200），不会被 `ImagePullBackOff` 卡住 |
| 18 | `secretKeyRef` / `Service.selector` / `serviceName` 引用 | 9 处 secret 引用、3 个 Service selector、1 个 StatefulSet serviceName 全部解析成功，无悬空引用 |
| 19 | 优雅关闭 | 服务端监听 `SIGINT/SIGTERM` 并 `srv.Shutdown`（5s 超时），小于默认 `terminationGracePeriodSeconds: 30`，滚动更新不会硬杀。无需额外配置 |
| 20 | NodePort `30082` 是否与常见端口冲突 | 不冲突。30082 落在默认 30000–32767 内，且不是任何知名服务的端口（知名的是 3000 Grafana、30080/30443 各类 ingress 默认、30000 常见 dashboard）。唯一真实风险是「同集群已被别的 Service 占用」和「节点防火墙未放通」，见第七节第 7 条 |
| 21 | `mongo:7.0.28` / `redis:7.0.11` 与 compose 一致 | 一致。注意 `single-host-deployment-fix.yml` 用的是 `mongo:7.0.4`（旧版兼容分支），k8s 用 7.0.28 更合理，两者数据目录兼容 |
| 22 | 端口映射 | `containerPort 8082` = `Service.port 8082` = `targetPort 8082` = 代码默认 `server.port 8082`，四者一致；`scan` 容器不监听端口（worker 模式），清单未写 `ports` 是正确的 |

---

## 附录：一致性自检

`all.yml` 与分文件是否漂移：

```bash
python3 - <<'PY'
import yaml, glob, sys
def key(d): return (d['kind'], d['metadata']['name'])
all_docs = {key(d): d for d in yaml.safe_load_all(open('k8s/all.yml')) if d}
bad = []
n = 0
for f in sorted(glob.glob('k8s/**/*.yml', recursive=True)):
    if f == 'k8s/all.yml': continue
    for d in yaml.safe_load_all(open(f)):
        if not d: continue
        n += 1
        if key(d) not in all_docs: bad.append((f, key(d), 'missing in all.yml'))
        elif all_docs[key(d)] != d: bad.append((f, key(d), 'content differs'))
print(f"checked {n} docs, all.yml has {len(all_docs)}")
print("DRIFT:", bad if bad else "none")
sys.exit(1 if bad else 0)
PY
```

期望输出：

```
checked 10 docs, all.yml has 10
DRIFT: none
```

CI 里建议串起来跑：

```bash
python3 -c "import yaml,glob; [list(yaml.safe_load_all(open(f))) for f in glob.glob('k8s/**/*.yml', recursive=True)]; print('yaml ok')" \
  && kubeconform -strict -summary k8s/all.yml
```
