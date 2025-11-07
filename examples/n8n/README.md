# n8n + xiao-world: Quick Setup

> **สำหรับคนที่รู้จัก n8n แล้ว** - เอกสารฉบับย่อ

เชื่อมต่อ n8n กับ xiao-world เพื่อเผยแพร่เนื้อหาจากเสี้ยวหงชูไปยัง Twitter, Facebook, TikTok, YouTube อัตโนมัติ

---

## Prerequisites

ต้องมี:
- ✅ Docker & Docker Compose
- ✅ **xiao-world รันอยู่ที่** `localhost:18060`
- ✅ API Keys สำหรับ platforms (Twitter, Facebook, etc.)

**เช็ค xiao-world:**
```bash
curl http://localhost:18060/mcp -X POST \
  -H "Content-Type: application/json" \
  -d '{"jsonrpc":"2.0","method":"initialize","params":{},"id":1}'
```

---

## Quick Setup (3 Steps)

### 1. รัน n8n

```bash
# ใช้ Docker Compose (แนะนำ)
cd examples/n8n
docker-compose up -d

# หรือใช้ Docker โดยตรง
docker run -d --name n8n \
  -p 5678:5678 \
  -v n8n_data:/home/node/.n8n \
  --add-host=host.docker.internal:host-gateway \
  n8nio/n8n:latest
```

**เข้าใช้งาน:** http://localhost:5678

**Setup ครั้งแรก:**
- กรอก Email, Password (จดไว้!)
- คลิก Continue

### 2. Import Workflow

1. เข้า n8n → **Workflows** tab
2. คลิก **"+ Add workflow"** → **"Import from file"**
3. เลือกไฟล์: `xiao-world-workflow.json`
4. คลิก **Import**
5. คลิก **Save**

### 3. ตั้งค่า Feed ID

1. คลิก node **"📝 ตั้งค่า Feed ID"**
2. แก้ไข `feed_id` และ `xsec_token`:
   - `feed_id`: รหัสโพสต์จากเสี้ยวหงชู
   - `xsec_token`: Token สำหรับ API
3. คลิก **Save**

**วิธีหา feed_id และ xsec_token:**
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

---

## Usage

**รัน workflow:**
1. คลิกปุ่ม **"Execute Workflow"** (⚡) ที่มุมบนขวา
2. รอ 8-10 วินาที
3. ดูผลลัพธ์ที่ node สุดท้าย

**ตรวจสอบผลลัพธ์:**
- เช็คที่ Twitter profile
- เช็คที่ Facebook page

---

## Workflows Available

### 1. `xiao-world-workflow.json` (Multi-Platform Publisher)

**8 Nodes:**
1. Manual Trigger - เริ่มต้น
2. Set Feed ID - ตั้งค่า feed_id และ xsec_token
3. Get Feed Detail - ดึงข้อมูลจากเสี้ยวหงชู (MCP API)
4. Parse Data - แปลงข้อมูล
5. Publish to Twitter - โพสต์ไป Twitter
6. Publish to Facebook - โพสต์ไป Facebook
7. Merge Results - รวมผลลัพธ์
8. Format Output - แสดงสรุป

**รองรับ Platforms:**
- Twitter ✅
- Facebook ✅
- TikTok (เพิ่ม node ได้)
- YouTube (เพิ่ม node ได้)

**ใช้เวลา:** ~8-10 วินาที

---

## Customization

### เพิ่ม Platform อื่น

เพิ่ม HTTP Request node:
- **URL:** `http://host.docker.internal:18060/mcp`
- **Method:** POST
- **Body:**
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

### ตั้งเวลาอัตโนมัติ

1. ลบ Manual Trigger node
2. เพิ่ม **Schedule Trigger** node
3. ตั้งค่า:
   - Daily: 10:00 AM
   - Timezone: Asia/Bangkok
4. เชื่อมต่อกับ node ถัดไป
5. คลิก **"Active"** ที่มุมบนขวา

---

## Troubleshooting

### ❌ Connection Refused

**สาเหตุ:** xiao-world ไม่รัน หรือ URL ผิด

**แก้:**
- เช็คว่า xiao-world รันอยู่
- Mac/Windows: ใช้ `host.docker.internal`
- Linux: ใช้ `172.17.0.1`

### ❌ Invalid feed_id

**สาเหตุ:** feed_id หรือ xsec_token ผิด/หมดอายุ

**แก้:**
- ใช้ `list_feeds` หา feed_id ใหม่
- ขอ xsec_token ใหม่

### ❌ Platform not enabled

**สาเหตุ:** ไม่ได้ตั้งค่า API keys

**แก้:**
- เช็คไฟล์ `.env` ของ xiao-world
- เพิ่ม API keys ให้ครบ

---

## Files

```
examples/n8n/
├── README.md                      # เอกสารนี้
├── docker-compose.yml             # Docker setup สำหรับ n8n
├── xiao-world-workflow.json       # Workflow พร้อมใช้งาน
└── images/                        # Screenshots (optional)
```

---

## Next Steps

- ปรับแต่ง workflow ตามต้องการ
- เพิ่ม AI translation node (OpenAI, Claude)
- ตั้ง schedule สำหรับโพสต์อัตโนมัติ
- Export workflow เพื่อ backup

---

## Links

- [n8n Documentation](https://docs.n8n.io)
- [xiao-world GitHub](https://github.com/huge8888/xiao-world)
- [MCP Protocol](https://modelcontextprotocol.io)

---

**พร้อมใช้งานแล้ว!** 🎉

Made with ❤️ for Thai Community 🇹🇭
