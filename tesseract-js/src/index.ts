import { Hono } from 'hono'
import z from 'zod'
import { createWorker, PSM } from 'tesseract.js'

const app = new Hono()

app.get('/status', (c) => {
  return c.json({
    status: 'Healthy',
    timestamp: new Date().toISOString(),
    version: "0.1.0"
  })
})

const UploadSchema = z.object({
  base64Image: z.string().regex(/^data:image\/([a-zA-Z]*);base64,[A-Za-z0-9+\/=]+$/),
  lang: z.string(),
})

app.post('/upload', async (c) => {
  const body = await c.req.json();

  const { data, error } = UploadSchema.safeParse(body)
  if (!data) {
    return c.json({ error: 'Invalid input', details: error.message }, 400)
  }

  const worker = await createWorker(data.lang)
  await worker.setParameters({
    tessedit_pageseg_mode: PSM.AUTO_OSD,
  });
  const ret = await worker.recognize(data.base64Image, {}, { blocks: true })
  worker.terminate()

  return c.json(ret.data);
})

export default app
