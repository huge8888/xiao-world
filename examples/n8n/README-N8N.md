# 🤖 Xiao-World + n8n: ระบบเผยแพร่อัตโนมัติแบบ No-Code

> **ไม่ต้องเขียนโค้ด! ลากวาง ตั้งค่า เสร็จ!**
>
> ใช้ n8n สร้าง workflow อัตโนมัติเพื่อดึงเนื้อหาจากเสี้ยวหงชู แปลภาษา และเผยแพร่ไปหลายแพลตฟอร์มพร้อมกัน

[![n8n](https://img.shields.io/badge/n8n-Workflow_Automation-blue?style=flat-square)](https://n8n.io)
[![Docker](https://img.shields.io/badge/Docker-Ready-2496ED?style=flat-square&logo=docker)](https://www.docker.com/)

---

## 📋 สารบัญ

- [n8n คืออะไร?](#n8n-คืออะไร)
- [ทำอะไรได้บ้าง?](#ทำอะไรได้บ้าง)
- [ติดตั้งอย่างไร?](#ติดตั้งอย่างไร)
- [ใช้งานอย่างไร?](#ใช้งานอย่างไร)
- [ตัวอย่าง Workflow](#ตัวอย่าง-workflow)
- [แก้ปัญหา](#แก้ปัญหา)
- [FAQ](#faq)

---

## 🤔 n8n คืออะไร?

**n8n** = เครื่องมือสร้าง workflow อัตโนมัติแบบ visual (คล้าย Zapier แต่ self-hosted ฟรี!)

### ทำไมต้องใช้ n8n?

| ข้อดี | คำอธิบาย |
|-------|----------|
| 🎨 **Visual Editor** | ลากวาง node ไม่ต้องเขียนโค้ด |
| 🆓 **ฟรี 100%** | Self-hosted ไม่มีค่าใช้จ่าย |
| 🔌 **Integration เยอะ** | เชื่อมต่อได้หลายร้อยบริการ |
| 🤖 **AI-Ready** | รองรับ AI APIs หลายตัว |
| 💾 **ควบคุมเองได้** | ข้อมูลอยู่ในเครื่องคุณ |

### เหมาะกับใคร?

- ✅ มือใหม่ที่ไม่เคยเขียนโค้ด
- ✅ Content Creator ที่อยากประหยัดเวลา
- ✅ นักการตลาดที่อยากทำงานอัตโนมัติ
- ✅ ธุรกิจเล็กที่งบจำกัด

---

## ✨ ทำอะไรได้บ้าง?

### 🎯 Workflow ที่เราเตรียมไว้ให้

**Workflow 1: เผยแพร่หลายแพลตฟอร์ม (แนะนำ!)**

```
1️⃣ ดึงโพสต์จากเสี้ยวหงชู
         ↓
2️⃣ แปลภาษาอัตโนมัติ (ถ้าต้องการ)
         ↓
3️⃣ เผยแพร่ไป Twitter, Facebook พร้อมกัน
         ↓
4️⃣ แสดงผลลัพธ์
```

**ใช้เวลาทั้งหมด:** ~10 วินาที ⚡

---

## ⚡ Quick Start (สำหรับผู้ที่รู้จัก n8n แล้ว)

**ถ้าคุณเคยใช้ n8n มาก่อน เริ่มได้เลยใน 3 ขั้นตอน:**

```bash
# 1. รัน n8n
cd examples/n8n && docker-compose up -d

# 2. เปิดเบราว์เซอร์
# http://localhost:5678

# 3. Import workflow
# Import file: xiao-world-workflow.json
# แก้ไข node "📝 ตั้งค่า Feed ID" → ใส่ feed_id และ xsec_token
# Execute! ⚡
```

**สำหรับมือใหม่:** อ่านคู่มือโดยละเอียดด้านล่าง ↓

---

## 🚀 ติดตั้งอย่างไร?

### ขั้นตอนที่ 1: ตรวจสอบความพร้อม

คุณต้องมี:

- ✅ **Docker & Docker Compose** (ติดตั้งจาก [docker.com](https://www.docker.com/get-started))
- ✅ **xiao-world** กำลังรันอยู่ที่ `localhost:18060`
- ✅ **API Keys** (Twitter, Facebook, etc.) ที่ตั้งค่าใน xiao-world แล้ว

**ตรวจสอบ xiao-world:**

```bash
curl http://localhost:18060/mcp -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"initialize","params":{},"id":1}'
```

ถ้าได้ response กลับมา แสดงว่าพร้อมแล้ว! ✅

---

### ขั้นตอนที่ 2: ติดตั้ง n8n

#### วิธีที่ 1: ใช้ Docker Compose (แนะนำ! ⭐)

```bash
# 1. เข้าไปในโฟลเดอร์ examples/n8n
cd examples/n8n

# 2. รัน n8n
docker-compose up -d

# 3. รอสักครู่ให้ n8n เริ่มต้น (~10 วินาที)
```

**เสร็จแล้ว!** 🎉

#### วิธีที่ 2: รันด้วย Docker โดยตรง

```bash
docker run -d \
  --name xiao-world-n8n \
  -p 5678:5678 \
  -v n8n_data:/home/node/.n8n \
  --add-host=host.docker.internal:host-gateway \
  n8nio/n8n:latest
```

---

### ขั้นตอนที่ 3: เปิด n8n ครั้งแรก (โดยละเอียด)

#### 3.1 เปิดเบราว์เซอร์

1. **เปิดเบราว์เซอร์** (Chrome, Firefox, Safari ก็ได้)
2. **พิมพ์ URL:**
   ```
   http://localhost:5678
   ```
3. **กด Enter**

#### 3.2 หน้าจอ Setup ครั้งแรก

**จะเห็นหน้าจอต้อนรับ n8n** พร้อมฟอร์มลงทะเบียน:

```
┌─────────────────────────────────────┐
│   Welcome to n8n!                   │
│                                     │
│   Owner account setup               │
│                                     │
│   Email: [____________]             │
│   First name: [____________]        │
│   Last name: [____________]         │
│   Password: [____________]          │
│                                     │
│   [X] Subscribe to newsletter       │
│                                     │
│        [ Continue ]                 │
└─────────────────────────────────────┘
```

**กรอกข้อมูล:**

1. **Email:** อีเมลของคุณ (ใช้อะไรก็ได้)
   - ตัวอย่าง: `myemail@gmail.com`
   - **จดไว้!** ใช้ login ครั้งต่อไป

2. **First name:** ชื่อจริง
   - ตัวอย่าง: `สมชาย`

3. **Last name:** นามสกุล
   - ตัวอย่าง: `ใจดี`

4. **Password:** รหัสผ่าน (ยาวอย่างน้อย 8 ตัว)
   - ตัวอย่าง: `MyP@ssw0rd123`
   - **จดไว้!** ใช้ login ครั้งต่อไป

5. **Newsletter:** ✅ หรือ ❌ ตามใจชอบ

6. **คลิก "Continue"**

#### 3.3 หน้าจอหลังจาก Setup

**จะเห็นหน้าจอ Dashboard ของ n8n:**

```
┌─────────────────────────────────────────────────────────┐
│  n8n                                   [Profile] [Help]  │
├─────────────────────────────────────────────────────────┤
│  Workflows    Credentials    Executions                 │
├─────────────────────────────────────────────────────────┤
│                                                          │
│     Welcome to n8n! 👋                                  │
│                                                          │
│     Get started with your first workflow:                │
│                                                          │
│     [ + Add workflow ]                                   │
│     [ Import from file ]                                 │
│     [ Use template ]                                     │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

**เมนูที่สำคัญ:**

- **Workflows** = รายการ workflow ทั้งหมด (ตอนนี้ยังว่างเปล่า)
- **Credentials** = เก็บ API keys ต่างๆ
- **Executions** = ประวัติการรัน workflow

**ตอนนี้พร้อมใช้งานแล้ว!** 🚀 เดี๋ยวเราจะ import workflow ในขั้นตอนถัดไป

---

## 🎮 ใช้งานอย่างไร?

### ขั้นตอนที่ 1: Import Workflow (โดยละเอียด)

#### 1.1 เข้าสู่หน้า Workflows

**ถ้าคุณอยู่ที่หน้า Dashboard:**

1. คลิกที่แท็บ **"Workflows"** ที่มุมบนซ้าย
2. จะเห็นหน้าจอว่างเปล่า (เพราะยังไม่มี workflow)

**ถ้าคุณปิดเบราว์เซอร์ไปแล้ว:**

1. เปิดเบราว์เซอร์ใหม่
2. ไปที่ `http://localhost:5678`
3. Login ด้วย Email และ Password ที่ตั้งไว้
4. คลิกแท็บ **"Workflows"**

#### 1.2 เริ่ม Import Workflow

**หน้าจอที่จะเห็น:**

```
┌─────────────────────────────────────────────────────────┐
│  n8n                                   [Profile] [Help]  │
├─────────────────────────────────────────────────────────┤
│  Workflows    Credentials    Executions                 │
├─────────────────────────────────────────────────────────┤
│                                                          │
│  My workflows                          [ + Add workflow ▼] │
│                                                          │
│  No workflows yet                                        │
│                                                          │
└─────────────────────────────────────────────────────────┘
```

**ทำตามขั้นตอน:**

1. **คลิกปุ่ม** `[ + Add workflow ▼]` ที่มุมบนขวา

2. **เมนู dropdown จะปรากฏ:**
   ```
   ┌─────────────────────────┐
   │  + Add workflow         │
   │  ─────────────────────  │
   │  ○ Create new           │
   │  ○ Import from file     │  ← คลิกตรงนี้!
   │  ○ Import from URL      │
   └─────────────────────────┘
   ```

3. **คลิกที่** `Import from file`

#### 1.3 เลือกไฟล์ Workflow

**หน้าต่างเลือกไฟล์จะเปิดขึ้น:**

```
┌────────────────────────────────────────┐
│  Select workflow file to import        │
│                                        │
│  Accepted formats: JSON                │
│                                        │
│  [ Browse files... ]                   │
│  or drag and drop file here            │
└────────────────────────────────────────┘
```

**วิธีเลือกไฟล์:**

**วิธีที่ 1: คลิก Browse**

1. คลิก **"Browse files..."**
2. ไปที่โฟลเดอร์ `xiao-world/examples/n8n/`
3. เลือกไฟล์ **`xiao-world-workflow.json`**
4. คลิก **"เปิด"** (หรือ "Open")

**วิธีที่ 2: ลากไฟล์มาวาง (Drag & Drop)**

1. เปิด File Explorer/Finder
2. ไปที่โฟลเดอร์ `xiao-world/examples/n8n/`
3. ลากไฟล์ `xiao-world-workflow.json`
4. วางลงในกล่อง "drag and drop file here"

#### 1.4 ตรวจสอบและ Import

**หลังจากเลือกไฟล์ จะเห็น:**

```
┌────────────────────────────────────────┐
│  Import workflow                       │
│                                        │
│  File: xiao-world-workflow.json        │
│                                        │
│  Workflow name:                        │
│  🌍 Xiao-World: เผยแพร่เนื้อหา...     │
│                                        │
│  This workflow contains 8 nodes        │
│                                        │
│       [ Cancel ]    [ Import ]         │
└────────────────────────────────────────┘
```

**คลิกปุ่ม** `[ Import ]`

#### 1.5 Workflow ถูก Import แล้ว!

**จะเห็นหน้าจอ Workflow Editor:**

```
┌────────────────────────────────────────────────────────────────┐
│  🌍 Xiao-World: เผยแพร่เนื้อหา...        [ Save ] [ Execute ]  │
├────────────────────────────────────────────────────────────────┤
│                                                                 │
│   [เริ่มต้น Manual] → [📝 ตั้งค่า Feed ID]                   │
│                              ↓                                  │
│                  [📥 ดึงข้อมูลจากเสี้ยวหงชู]                  │
│                              ↓                                  │
│                      [🔧 แปลงข้อมูล]                           │
│                     ↙              ↘                           │
│      [🐦 โพสต์ไป Twitter]    [📘 โพสต์ไป Facebook]           │
│                     ↘              ↙                           │
│                    [📊 รวมผลลัพธ์]                             │
│                              ↓                                  │
│                    [✅ สรุปผลลัพธ์]                            │
│                                                                 │
└────────────────────────────────────────────────────────────────┘
```

**สิ่งที่ต้องทำ:**

1. **คลิก "Save"** เพื่อบันทึก workflow (ตรงมุมบนขวา)
2. จะมี popup ขึ้นมา:
   ```
   ┌──────────────────────┐
   │  Workflow saved! ✓   │
   └──────────────────────┘
   ```

**เสร็จแล้ว!** Workflow พร้อมใช้งาน! 🎉

**Tip:** ดู node แต่ละตัวด้วยการคลิกที่มันเพื่อดูรายละเอียด

---

### ขั้นตอนที่ 2: ตรวจสอบการเชื่อมต่อ (โดยละเอียด)

#### 2.1 ทำความเข้าใจ Node

**ก่อนอื่น มาดู node ที่สำคัญกันก่อน:**

```
[📥 ดึงข้อมูลจากเสี้ยวหงชู] = Node ที่เรียก MCP API ของ xiao-world
```

Node นี้จะเชื่อมต่อไปที่ xiao-world server ที่ `localhost:18060`

#### 2.2 เปิด Node Editor

**ทำตามขั้นตอน:**

1. **คลิกที่ node** `📥 ดึงข้อมูลจากเสี้ยวหงชู` (คลิกครั้งเดียว)

2. **Panel ด้านขวาจะเปิดขึ้น:**

```
┌──────────────────────────────────────────────┐
│  📥 ดึงข้อมูลจากเสี้ยวหงชู                  │
│  ────────────────────────────────────────   │
│                                              │
│  Node Type: HTTP Request                     │
│                                              │
│  Parameters:                                 │
│  ┌────────────────────────────────────────┐ │
│  │ Method: POST                    ▼      │ │
│  ├────────────────────────────────────────┤ │
│  │ URL:                                   │ │
│  │ http://host.docker.internal:18060/mcp  │ │
│  ├────────────────────────────────────────┤ │
│  │ Authentication: None            ▼      │ │
│  ├────────────────────────────────────────┤ │
│  │ Send Headers: [X]                      │ │
│  │   Content-Type: application/json       │ │
│  ├────────────────────────────────────────┤ │
│  │ Send Body: [X]                         │ │
│  │   Body Format: JSON                    │ │
│  │   {...}                                │ │
│  └────────────────────────────────────────┘ │
│                                              │
│              [ Execute Node ]                │
└──────────────────────────────────────────────┘
```

#### 2.3 ตรวจสอบ URL

**สิ่งที่ต้องเช็ค:**

✅ **ถ้าใช้ Docker Compose (แนะนำ):**
```
URL: http://host.docker.internal:18060/mcp
```
- `host.docker.internal` = เชื่อมต่อไปที่ host machine จาก Docker
- `18060` = port ของ xiao-world (ค่า default)
- `/mcp` = MCP endpoint

✅ **ถ้า xiao-world รันบน Linux:**
```
URL: http://172.17.0.1:18060/mcp
```
- `172.17.0.1` = Docker bridge IP บน Linux

✅ **ถ้า xiao-world รันบน server อื่น:**
```
URL: http://192.168.1.100:18060/mcp
```
- แทนที่ `192.168.1.100` ด้วย IP จริงของ server

#### 2.4 ทดสอบการเชื่อมต่อ

**ทดสอบ node เดี่ยวๆ:**

1. คลิกที่ node `📥 ดึงข้อมูลจากเสี้ยวหงชู`
2. ที่ panel ด้านขวา เลื่อนลงล่างสุด
3. คลิกปุ่ม **"Execute Node"**

**ถ้าเชื่อมต่อสำเร็จ:**

```
┌──────────────────────────────────────────┐
│  ✓ Node executed successfully            │
│                                          │
│  Output 1:                               │
│  {                                       │
│    "jsonrpc": "2.0",                     │
│    "result": {                           │
│      "content": [...]                    │
│    },                                    │
│    "id": 1                               │
│  }                                       │
└──────────────────────────────────────────┘
```

**ถ้าเชื่อมต่อไม่สำเร็จ:**

```
┌──────────────────────────────────────────┐
│  ✗ Execution failed                      │
│                                          │
│  Error: connect ECONNREFUSED             │
│  127.0.0.1:18060                         │
│                                          │
│  Possible causes:                        │
│  • xiao-world is not running             │
│  • Wrong URL or port                     │
│  • Firewall blocking connection          │
└──────────────────────────────────────────┘
```

**วิธีแก้:**
- ดูที่ส่วน [แก้ปัญหา](#แก้ปัญหา) ด้านล่าง
- เช็คว่า xiao-world รันอยู่: `curl http://localhost:18060/mcp`

#### 2.5 ตรวจสอบ JSON Body

**เช็คว่า JSON ถูกต้อง:**

1. คลิกที่ **"Body"** ใน node
2. เลื่อนลงมาดู **"JSON"** field
3. ควรเห็น:

```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "get_feed_detail",
    "arguments": {
      "feed_id": "{{ $json.feed_id }}",
      "xsec_token": "{{ $json.xsec_token }}"
    }
  },
  "id": 1
}
```

**คำอธิบาย:**
- `{{ $json.feed_id }}` = ดึงค่าจาก node ก่อนหน้า
- `{{ $json.xsec_token }}` = ดึงค่าจาก node ก่อนหน้า
- **ห้ามแก้ไข!** มันจะดึงค่าอัตโนมัติ

**เสร็จแล้ว!** การเชื่อมต่อพร้อมใช้งาน ✅

---

### ขั้นตอนที่ 3: ใส่ Feed ID (โดยละเอียด)

#### 3.1 ทำความเข้าใจ Feed ID

**Feed ID คืออะไร?**
- รหัสประจำตัวของโพสต์จากเสี้ยวหงชู
- ใช้เพื่อระบุโพสต์ที่ต้องการดึงมา

**xsec_token คืออะไร?**
- Token สำหรับยืนยันตัวตนกับ API ของเสี้ยวหงชู
- ใช้ร่วมกับ feed_id ทุกครั้ง

#### 3.2 วิธีหา Feed ID และ xsec_token

**วิธีที่ 1: ใช้ MCP Client (เช่น Claude Desktop)**

1. เปิด Claude Desktop
2. พิมพ์คำสั่ง:
   ```
   ใช้ tool list_feeds เพื่อแสดงโพสต์ล่าสุด
   ```
3. จะได้ผลลัพธ์:
   ```json
   [
     {
       "feed_id": "6751234567890abcdef",
       "xsec_token": "XYZ123abc...",
       "title": "หัวข้อโพสต์",
       "time": "2025-01-07"
     }
   ]
   ```
4. **คัดลอก** `feed_id` และ `xsec_token`

**วิธีที่ 2: ใช้ curl (สำหรับผู้ชำนาญ)**

```bash
curl http://localhost:18060/mcp -X POST \
  -H "Content-Type: application/json" \
  -d '{
    "jsonrpc": "2.0",
    "method": "tools/call",
    "params": {
      "name": "list_feeds",
      "arguments": {}
    },
    "id": 1
  }'
```

**ตัวอย่างผลลัพธ์:**
```json
{
  "feed_id": "6751234567890abcdef",
  "xsec_token": "XYZ123abc..."
}
```

#### 3.3 เปิด Node "📝 ตั้งค่า Feed ID"

**ทำตามขั้นตอน:**

1. **กลับไปที่หน้า workflow editor**
2. **คลิกที่ node** `📝 ตั้งค่า Feed ID` (node สีเหลืองที่ 2)

**Panel ด้านขวาจะเปิดขึ้น:**

```
┌──────────────────────────────────────────────┐
│  📝 ตั้งค่า Feed ID                          │
│  ────────────────────────────────────────   │
│                                              │
│  Node Type: Set                              │
│                                              │
│  Assignments:                                │
│  ┌────────────────────────────────────────┐ │
│  │ Name: feed_id                          │ │
│  │ Type: String                           │ │
│  │ Value: 6751234567890abcdef             │ │
│  │                              [Edit]    │ │
│  ├────────────────────────────────────────┤ │
│  │ Name: xsec_token                       │ │
│  │ Type: String                           │ │
│  │ Value: your_token_here                 │ │
│  │                              [Edit]    │ │
│  └────────────────────────────────────────┘ │
│                                              │
│          [ + Add Assignment ]                │
└──────────────────────────────────────────────┘
```

#### 3.4 แก้ไข Feed ID

**แก้ไขค่า feed_id:**

1. **คลิกที่** assignment `feed_id`
2. **คลิกปุ่ม [Edit]** ทางขวา
3. **จะเห็น popup:**

```
┌────────────────────────────────────┐
│  Edit assignment                   │
│                                    │
│  Name: feed_id                     │
│                                    │
│  Type: String              ▼       │
│                                    │
│  Value:                            │
│  ┌──────────────────────────────┐ │
│  │ 6751234567890abcdef          │ │
│  └──────────────────────────────┘ │
│                                    │
│     [ Cancel ]    [ Update ]       │
└────────────────────────────────────┘
```

4. **ลบค่าเก่า** `6751234567890abcdef`
5. **วาง feed_id ของคุณ** (ที่คัดลอกมา)
   - ตัวอย่าง: `675a8b9c01d2e3f4567890ab`
6. **คลิก [Update]**

#### 3.5 แก้ไข xsec_token

**แก้ไขค่า xsec_token:**

1. **คลิกที่** assignment `xsec_token`
2. **คลิกปุ่ม [Edit]**
3. **ลบค่าเก่า** `your_token_here`
4. **วาง xsec_token ของคุณ**
   - ตัวอย่าง: `XYZ789def456abc123...` (จะยาวมาก!)
5. **คลิก [Update]**

**ตัวอย่างหลังแก้ไขแล้ว:**

```
┌────────────────────────────────────────────┐
│ Assignments:                               │
│ ┌────────────────────────────────────────┐ │
│ │ Name: feed_id                          │ │
│ │ Value: 675a8b9c01d2e3f4567890ab        │ │ ✅
│ ├────────────────────────────────────────┤ │
│ │ Name: xsec_token                       │ │
│ │ Value: XYZ789def456abc123...           │ │ ✅
│ └────────────────────────────────────────┘ │
└────────────────────────────────────────────┘
```

#### 3.6 บันทึกการแก้ไข

**อย่าลืมบันทึก!**

1. **คลิกปุ่ม "Save"** ที่มุมบนขวาของหน้า workflow
2. จะมี notification:
   ```
   ┌──────────────────────┐
   │  Workflow saved! ✓   │
   └──────────────────────┘
   ```

**เสร็จแล้ว!** Feed ID ถูกตั้งค่าเรียบร้อย ✅

**Tips:**
- ถ้าต้องการเปลี่ยน feed_id ครั้งต่อไป แค่มาแก้ไขที่ node นี้อีกครั้ง
- xsec_token อาจหมดอายุ ต้องขอ token ใหม่จาก `list_feeds`

---

### ขั้นตอนที่ 4: รัน Workflow! (โดยละเอียด)

#### 4.1 เตรียมตัวก่อนรัน

**เช็คให้แน่ใจว่า:**

✅ xiao-world รันอยู่ที่ `localhost:18060`
✅ Feed ID และ xsec_token ถูกต้อง
✅ Workflow ถูก save แล้ว
✅ ไม่มี error ที่ node ใดๆ (ไม่มีเครื่องหมายแดง)

#### 4.2 คลิกปุ่ม Execute

**ทำตามขั้นตอน:**

1. **กลับไปที่หน้า workflow editor**
2. **คลิกปุ่ม "Execute Workflow"** ที่มุมบนขวา

```
┌────────────────────────────────────────────┐
│  🌍 Xiao-World...    [ Save ] [ Execute ⚡] │ ← คลิกตรงนี้!
└────────────────────────────────────────────┘
```

**จะเห็น popup ยืนยัน:**

```
┌────────────────────────────────────┐
│  Execute workflow?                 │
│                                    │
│  This will run the workflow with   │
│  the current settings.             │
│                                    │
│     [ Cancel ]    [ Execute ]      │
└────────────────────────────────────┘
```

3. **คลิก [Execute]**

#### 4.3 ระหว่างการทำงาน

**จะเห็น animation ที่แต่ละ node:**

```
┌────────────────────────────────────────────────────┐
│                                                     │
│  [เริ่มต้น Manual] ✓                               │
│         ↓                                           │
│  [📝 ตั้งค่า Feed ID] ⟳ (กำลังทำงาน...)          │
│         ↓                                           │
│  [📥 ดึงข้อมูล...] ⏸ (รอ...)                      │
│                                                     │
└────────────────────────────────────────────────────┘
```

**สัญลักษณ์:**
- ✓ = เสร็จแล้ว (สีเขียว)
- ⟳ = กำลังทำงาน (สีน้ำเงิน, มีแอนิเมชั่น)
- ⏸ = รอ node ก่อนหน้าเสร็จ (สีเทา)
- ✗ = เกิด error (สีแดง)

**ใช้เวลา:**
- ~2-3 วินาที ดึงข้อมูลจากเสี้ยวหงชู
- ~1 วินาที แปลงข้อมูล
- ~2-3 วินาที โพสต์แต่ละแพลตฟอร์ม
- **รวม: ~8-10 วินาที** ⚡

#### 4.4 ดูผลลัพธ์แต่ละ Node

**หลังจาก execute เสร็จ จะเห็น node ทุกตัวมี checkmark ✓**

**คลิกดู node แต่ละตัว:**

**Node 1: เริ่มต้น Manual**
```
┌──────────────────────────────────┐
│  ✓ Executed successfully         │
│                                  │
│  Output: 1 item                  │
│  { trigger: "manual" }           │
└──────────────────────────────────┘
```

**Node 2: 📝 ตั้งค่า Feed ID**
```
┌──────────────────────────────────┐
│  ✓ Executed successfully         │
│                                  │
│  Output: 1 item                  │
│  {                               │
│    "feed_id": "675a8b9c...",     │
│    "xsec_token": "XYZ789..."     │
│  }                               │
└──────────────────────────────────┘
```

**Node 3: 📥 ดึงข้อมูลจากเสี้ยวหงชู**
```
┌──────────────────────────────────────────┐
│  ✓ Executed successfully                 │
│                                          │
│  Output: 1 item                          │
│  {                                       │
│    "jsonrpc": "2.0",                     │
│    "result": {                           │
│      "content": [                        │
│        {                                 │
│          "type": "text",                 │
│          "text": "{                      │
│            \"title\": \"วิธีทำอาหาร\",  │
│            \"desc\": \"...\",            │
│            \"images\": [...]             │
│          }"                              │
│        }                                 │
│      ]                                   │
│    },                                    │
│    "id": 1                               │
│  }                                       │
└──────────────────────────────────────────┘
```

**Node 4: 🔧 แปลงข้อมูล**
```
┌──────────────────────────────────────────┐
│  ✓ Executed successfully                 │
│                                          │
│  Output: 1 item                          │
│  {                                       │
│    "feed_id": "675a8b9c...",             │
│    "title": "วิธีทำอาหาร",              │
│    "content": "วันนี้จะมาสอน...",       │
│    "images": [                           │
│      "https://...",                      │
│      "https://..."                       │
│    ],                                    │
│    "tags": ["อาหาร", "สูตรอาหาร"]       │
│  }                                       │
└──────────────────────────────────────────┘
```

**Node 5 & 6: 🐦 Twitter + 📘 Facebook**
```
┌──────────────────────────────────────────┐
│  ✓ Executed successfully (Twitter)       │
│                                          │
│  {                                       │
│    "result": {                           │
│      "content": [                        │
│        {                                 │
│          "text": "{                      │
│            \"platform\": \"twitter\",    │
│            \"status\": \"success\",      │
│            \"tweet_id\": \"123...\"      │
│          }"                              │
│        }                                 │
│      ]                                   │
│    }                                     │
│  }                                       │
└──────────────────────────────────────────┘
```

**Node 8: ✅ สรุปผลลัพธ์**
```
┌──────────────────────────────────────────┐
│  ✓ Executed successfully                 │
│                                          │
│  Output: 1 item                          │
│  {                                       │
│    "feed_info": {                        │
│      "feed_id": "675a8b9c...",           │
│      "title": "วิธีทำอาหาร",            │
│      "content": "วันนี้จะมาสอน..."      │
│    },                                    │
│    "publish_results": {                  │
│      "twitter": {                        │
│        "status": "success",              │
│        "tweet_id": "123..."              │
│      },                                  │
│      "facebook": {                       │
│        "status": "success",              │
│        "post_id": "456..."               │
│      }                                   │
│    },                                    │
│    "timestamp": "2025-01-07T10:30:00Z"   │
│  }                                       │
└──────────────────────────────────────────┘
```

#### 4.5 ตรวจสอบผลลัพธ์

**เช็คว่าโพสต์สำเร็จหรือไม่:**

1. **ดูที่ node สุดท้าย** (✅ สรุปผลลัพธ์)
2. **เช็ค status:**
   - `"status": "success"` ✅ = สำเร็จ
   - `"status": "error"` ❌ = ล้มเหลว

3. **ตรวจสอบที่ Platform จริง:**
   - Twitter: ไปดูที่ profile Twitter ของคุณ
   - Facebook: ไปดูที่ Facebook page ของคุณ

#### 4.6 ถ้าเกิด Error

**ถ้าเห็นสัญลักษณ์ ✗ (สีแดง) ที่ node ใดๆ:**

1. **คลิกที่ node นั้น**
2. **อ่าน error message:**

```
┌──────────────────────────────────────────┐
│  ✗ Execution failed                      │
│                                          │
│  Error: Invalid feed_id                  │
│                                          │
│  The feed_id provided does not exist     │
│  or has been deleted.                    │
│                                          │
│  Possible solutions:                     │
│  • Check if feed_id is correct           │
│  • Use list_feeds to get valid IDs       │
│  • Check if xsec_token is expired        │
└──────────────────────────────────────────┘
```

3. **แก้ไขตาม error message**
4. **รัน workflow ใหม่อีกครั้ง**

**Error ที่พบบ่อย:**
- `Invalid feed_id` → ใช้ `list_feeds` หา ID ใหม่
- `Invalid token` → ขอ xsec_token ใหม่
- `Platform not enabled` → เช็คไฟล์ `.env` ของ xiao-world
- `Connection refused` → เช็คว่า xiao-world รันอยู่หรือไม่

**ดูเพิ่มเติมที่:** [แก้ปัญหา](#แก้ปัญหา)

#### 4.7 ดูประวัติการรัน

**ต้องการดูประวัติ execution:**

1. คลิกแท็บ **"Executions"** ที่มุมบนซ้าย
2. จะเห็นรายการ execution ทั้งหมด:

```
┌────────────────────────────────────────────────────┐
│  Executions                                        │
├────────────────────────────────────────────────────┤
│  Workflow: 🌍 Xiao-World                          │
│                                                    │
│  ✓ 2025-01-07 10:30  Success   8 nodes   10s      │
│  ✓ 2025-01-07 09:15  Success   8 nodes   9s       │
│  ✗ 2025-01-07 09:10  Failed    3 nodes   2s       │
│  ✓ 2025-01-07 08:45  Success   8 nodes   11s      │
└────────────────────────────────────────────────────┘
```

3. **คลิกที่ execution ใดๆ** เพื่อดูรายละเอียด

**เสร็จสิ้น!** คุณรัน workflow สำเร็จแล้ว! 🎉

---

### ขั้นตอนที่ 5: ทำอะไรต่อหลังรันสำเร็จ?

#### 5.1 ตรวจสอบผลลัพธ์บน Platform จริง

**Twitter:**
1. เปิด https://twitter.com
2. Login เข้าบัญชีของคุณ
3. ไปที่ Profile → Tweets
4. จะเห็นโพสต์ที่เพิ่งโพสต์ไป!

**Facebook:**
1. เปิด https://facebook.com
2. ไปที่ Page ที่คุณตั้งค่าไว้
3. จะเห็นโพสต์ใหม่!

#### 5.2 ปรับแต่ง Workflow (ถ้าต้องการ)

**เพิ่ม Platform อื่น:**

1. **คลิกที่พื้นที่ว่าง** ใน workflow editor
2. **กด "+"** เพื่อเพิ่ม node
3. **ค้นหา** "HTTP Request"
4. **เพิ่ม node** และเชื่อมต่อ
5. **ตั้งค่า:**
   - URL: `http://host.docker.internal:18060/mcp`
   - Method: `POST`
   - Body: JSON-RPC request สำหรับ platform อื่น (เช่น TikTok, YouTube)

**ตัวอย่าง JSON สำหรับ TikTok:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "publish_to_tiktok",
    "arguments": {
      "feed_id": "{{ $json.feed_id }}",
      "xsec_token": "{{ $json.xsec_token }}"
    }
  },
  "id": 4
}
```

**ตัวอย่าง JSON สำหรับ YouTube:**
```json
{
  "jsonrpc": "2.0",
  "method": "tools/call",
  "params": {
    "name": "publish_to_youtube",
    "arguments": {
      "feed_id": "{{ $json.feed_id }}",
      "xsec_token": "{{ $json.xsec_token }}"
    }
  },
  "id": 5
}
```

#### 5.3 ตั้งเวลาให้รันอัตโนมัติ

**เปลี่ยนจาก Manual Trigger เป็น Schedule:**

1. **ลบ node** "เริ่มต้น Manual"
2. **เพิ่ม node** "Schedule Trigger":
   - กด "+" → ค้นหา "Schedule"
   - เลือก "Schedule Trigger"
3. **ตั้งค่าเวลา:**

```
┌────────────────────────────────────────┐
│  Schedule Trigger                      │
│                                        │
│  Trigger Interval: Daily       ▼       │
│                                        │
│  Trigger Time: 10:00 AM                │
│                                        │
│  Timezone: Asia/Bangkok        ▼       │
└────────────────────────────────────────┘
```

**ตัวอย่างการตั้งเวลา:**
- **ทุกวัน 10:00 น.** → `Daily`, `10:00 AM`
- **ทุก 6 ชั่วโมง** → `Hours`, `6`
- **จันทร์-ศุกร์ เวลา 09:00** → `Cron`, `0 9 * * 1-5`

4. **เชื่อมต่อ** Schedule → 📝 ตั้งค่า Feed ID
5. **Save และ Activate workflow**

**สำคัญ:** ถ้าใช้ schedule ต้องคลิกปุ่ม **"Active"** ที่มุมบนขวา!

```
┌────────────────────────────────────────┐
│  🌍 Xiao-World   [Active ●]  [Save]    │
└────────────────────────────────────────┘
```

#### 5.4 เพิ่ม AI Translation

**แปลภาษาก่อนโพสต์:**

1. **เพิ่ม node** "OpenAI" หลังจาก "🔧 แปลงข้อมูล"
2. **ตั้งค่า:**
   - Model: `gpt-3.5-turbo`
   - Message: "แปลข้อความนี้เป็นภาษาไทย: {{ $json.content }}"
3. **เชื่อมต่อ:** แปลงข้อมูล → AI Translation → โพสต์

**ต้องมี:** OpenAI API Key (ใส่ใน Credentials)

#### 5.5 บันทึก Workflow

**Export เพื่อ backup:**

1. คลิก **"..." (More options)** ที่มุมบนขวา
2. เลือก **"Download"**
3. บันทึกไฟล์ JSON

**สามารถเอา workflow ไป import บนเครื่องอื่นได้!**

#### 5.6 Tips & Tricks

**เร็วขึ้น:**
- ใช้ "Execute Node" เพื่อทดสอบ node เดี่ยวๆ ก่อน
- ใช้ keyboard shortcut: `Ctrl+S` (Save), `Ctrl+Enter` (Execute)

**ปลอดภัยขึ้น:**
- เปิด Basic Auth ใน docker-compose.yml
- ใช้ HTTPS (ต้อง reverse proxy)
- เก็บ API keys ใน Credentials (ไม่ hardcode)

**ประหยัดเงิน:**
- ใช้ AI แค่เมื่อจำเป็น
- Cache ผลลัพธ์ที่ซ้ำ
- ตั้งค่า retry limit

---

## 🎨 ตัวอย่าง Workflow

### Workflow ที่ 1: เผยแพร่พร้อมกัน 2 แพลตฟอร์ม

**สิ่งที่ทำ:**

```
ดึงโพสต์จากเสี้ยวหงชู
    ↓
แปลง/ประมวลผลข้อมูล
    ↓
โพสต์ไป Twitter + Facebook พร้อมกัน
    ↓
แสดงผลลัพธ์
```

**การทำงาน:**

1. **Manual Trigger** → กดปุ่มเพื่อเริ่ม
2. **Set Feed ID** → ตั้งค่า Feed ID และ Token
3. **Get Feed Detail** → เรียก MCP API ดึงข้อมูล
4. **Parse Data** → แปลงข้อมูลให้ใช้งานง่าย
5. **Publish** → เผยแพร่ไป 2 แพลตฟอร์มพร้อมกัน
6. **Merge & Format** → รวมผลลัพธ์และแสดงผล

---

### Workflow ที่ 2: กำหนดเวลาโพสต์ (Coming Soon!)

**สิ่งที่ทำ:**

```
ตั้งเวลา (เช่น ทุกวันเวลา 10:00 น.)
    ↓
ค้นหาโพสต์ยอดนิยม
    ↓
เลือกโพสต์ที่ดีที่สุด
    ↓
เผยแพร่อัตโนมัติ
```

---

## 🐛 แก้ปัญหา

### ❌ ปัญหา: เชื่อมต่อ xiao-world ไม่ได้

**อาการ:**
```
Error: connect ECONNREFUSED
```

**วิธีแก้:**

1. **เช็คว่า xiao-world รันอยู่:**
   ```bash
   curl http://localhost:18060/mcp
   ```

2. **ถ้าใช้ Docker Desktop (Mac/Windows):**
   - ใช้ `http://host.docker.internal:18060/mcp` ✅

3. **ถ้าใช้ Linux:**
   - หาIP ของ host:
     ```bash
     ip addr show docker0 | grep inet
     ```
   - ใช้ `http://172.17.0.1:18060/mcp`

4. **ถ้ารันบน server อื่น:**
   - เปลี่ยนเป็น IP จริง เช่น `http://192.168.1.100:18060/mcp`

---

### ❌ ปัญหา: Import Workflow ไม่ได้

**วิธีแก้:**

1. ดาวน์โหลดไฟล์ `xiao-world-workflow.json` ใหม่
2. เช็คว่าไฟล์เป็น JSON ที่ถูกต้อง
3. ลอง copy-paste เนื้อหาไฟล์ไปที่ n8n โดยตรง:
   - Workflow → Import from URL → Paste JSON

---

### ❌ ปัญหา: Workflow รันไม่สำเร็จ

**วิธีแก้:**

1. **เช็คแต่ละ node:**
   - คลิกที่ node ที่ error
   - อ่าน error message
   - แก้ไขตาม message

2. **Error ที่พบบ่อย:**

| Error | สาเหตุ | วิธีแก้ |
|-------|--------|---------|
| `Invalid feed_id` | Feed ID ผิด | ใช้ `list_feeds` เพื่อหา ID ที่ถูกต้อง |
| `Invalid token` | xsec_token หมดอายุ | ขอ token ใหม่จาก `list_feeds` |
| `Platform not enabled` | ไม่ได้เปิดใช้แพลตฟอร์ม | เช็คไฟล์ `.env` ว่าตั้งค่าครบหรือยัง |

3. **ยังไม่ได้:**
   - ลองรัน xiao-world ด้วย MCP client อื่น (Claude Desktop)
   - เช็คว่าทำงานปกติหรือไม่

---

## ❓ FAQ

### Q1: n8n ใช้ฟรีได้จริงหรือ?

**ตอบ:** ใช้ฟรี 100%! เพราะเป็น self-hosted (รันในเครื่องคุณเอง)

**ต้องเสียเงินตรงไหน:**
- ❌ n8n: **ฟรี!**
- ❌ xiao-world: **ฟรี!**
- ✅ AI APIs: **เสียเงิน** ($0.15-0.25 ต่อ 1,000 คำ) ถ้าใช้
- ✅ Social media APIs: **ฟรีหมด!**

---

### Q2: ต้องเปิด n8n ทิ้งไว้ตลอดเวลาหรือ?

**ตอบ:**
- **ถ้าใช้ manual trigger** → เปิดเฉพาะตอนใช้งาน
- **ถ้าใช้ schedule (ตั้งเวลา)** → ต้องเปิดทิ้งไว้

**Tips:** รันบน VPS หรือ Server จะดีกว่า (เปิดได้ 24/7)

---

### Q3: ใช้ร่วมกับ AI ได้หรือไม่?

**ตอบ:** ได้! n8n มี nodes สำหรับ AI หลายตัว:

- ✅ OpenAI (ChatGPT)
- ✅ Anthropic (Claude)
- ✅ Google (Gemini)
- ✅ และอื่นๆ อีกมากมาย

**ตัวอย่าง:**
- ใช้ AI สร้าง caption ใหม่
- แปลภาษาก่อนโพสต์
- สร้าง hashtags อัตโนมัติ

---

### Q4: สร้าง Workflow เองได้หรือไม่?

**ตอบ:** ได้แน่นอน! n8n มี visual editor ที่ใช้งานง่าย

**ขั้นตอน:**
1. คลิก **"+ Add node"**
2. เลือก node type ที่ต้องการ
3. ลาก-วาง-เชื่อมต่อ
4. บันทึกและทดสอบ

**ไอเดีย Workflow:**
- ดึงโพสต์ยอดนิยมทุกวัน
- แปลภาษาอัตโนมัติ
- โพสต์ไปหลายแพลตฟอร์มตามเวลา
- ตรวจสอบ engagement และรายงาน

---

### Q5: ปลอดภัยหรือไม่?

**ตอบ:** ปลอดภัย!

- ✅ **ข้อมูลอยู่ในเครื่องคุณ** (ไม่ส่งไปไหน)
- ✅ **Open Source** (ตรวจสอบโค้ดได้เอง)
- ✅ **ควบคุมเองได้ทั้งหมด**

**เพิ่มความปลอดภัย:**
- เปิด Basic Auth (username/password)
- ใช้ HTTPS
- ตั้งค่า Firewall

---

### Q6: ใช้งานบน mobile ได้หรือไม่?

**ตอบ:** ได้! n8n เป็น web-based

- ✅ เปิดผ่านเบราว์เซอร์มือถือ
- ✅ กดปุ่ม execute ได้
- ⚠️ แต่แก้ไข workflow บน desktop จะสะดวกกว่า

---

### Q7: จะ backup workflow ยังไง?

**ตอบ:** มี 3 วิธี:

**วิธีที่ 1: Export JSON**
1. Workflow → ... → Download
2. บันทึกไฟล์ `.json`

**วิธีที่ 2: ใช้ Docker Volume**
```bash
docker run --rm -v n8n_data:/data -v $(pwd):/backup \
  busybox tar czf /backup/n8n-backup.tar.gz /data
```

**วิธีที่ 3: เก็บใน Git**
```bash
# เก็บ workflow ไว้ใน repo
git add examples/n8n/workflows/
git commit -m "backup workflows"
git push
```

---

### Q8: ช้าหรือเร็ว?

**ตอบ:** **เร็วมาก!** ⚡

**Benchmark:**
- เรียก MCP API: ~1-2 วินาที
- แปลงข้อมูล: <0.1 วินาที
- โพสต์แพลตฟอร์ม: ~2-3 วินาที แต่ละแพลตฟอร์ม

**รวม:** ~5-10 วินาทีต่อ workflow

---

### Q9: มี community หรือไม่?

**ตอบ:** มี!

**n8n Community:**
- 🌐 Forum: https://community.n8n.io
- 💬 Discord: มี community active มาก
- 📚 Docs: https://docs.n8n.io

**xiao-world Community:**
- 🐛 GitHub Issues: https://github.com/huge8888/xiao-world/issues
- 📧 Email: [ใส่ email ของคุณ]

---

## 📚 เอกสารเพิ่มเติม

### 🔗 Links ที่มีประโยชน์

- [n8n Documentation](https://docs.n8n.io)
- [n8n Community](https://community.n8n.io)
- [xiao-world README](../../README.md)
- [MCP Protocol](https://modelcontextprotocol.io)

### 📹 วิดีโอสอน (Coming Soon!)

- การติดตั้ง n8n + xiao-world
- สร้าง Workflow ขั้นพื้นฐาน
- สร้าง Workflow ขั้นสูง
- Tips & Tricks

---

## 🎁 Workflow Templates พิเศษ

เรากำลังพัฒนา workflow templates เพิ่มเติม:

- [ ] **Auto Daily Post** - โพสต์อัตโนมัติทุกวัน
- [ ] **Trending Monitor** - ติดตามโพสต์ยอดนิยม
- [ ] **Multi-language** - แปลและโพสต์หลายภาษา
- [ ] **Analytics Reporter** - รายงานสถิติอัตโนมัติ
- [ ] **Content Scheduler** - จัดคิวเนื้อหาล่วงหน้า

**ติดตามได้ที่:** [GitHub Issues](https://github.com/huge8888/xiao-world/issues)

---

## 💝 ขอบคุณ

- [n8n.io](https://n8n.io) - Workflow automation platform ที่เจ๋งมาก
- [Model Context Protocol](https://modelcontextprotocol.io) - MCP standard
- คุณทุกคน! ที่ใช้งาน xiao-world ❤️

---

<div align="center">

## 🌟 ให้ดาว xiao-world ถ้าชอบนะครับ!

**สนุกกับการทำงานอัตโนมัติ! 🤖✨**

Made with ❤️ for Thai Community 🇹🇭

</div>
