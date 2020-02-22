import http from 'http'
import request from 'supertest'
import { Connection } from 'typeorm'
import { getDb } from '../../database'
import { start, stop } from '../../support/server'
import { NuLinkNode, createNuLinkNode } from '../../entity/NuLinkNode'
import { JobRun } from '../../entity/JobRun'
import { TaskRun } from '../../entity/TaskRun'
import { createJobRun } from '../../factories'

let server: http.Server
let db: Connection

beforeAll(async () => {
  db = await getDb()
  server = await start()
})
afterAll(done => stop(server, done))

describe('#index', () => {
  describe('with no runs', () => {
    it('returns empty', async () => {
      const response = await request(server).get('/api/v1/job_runs')
      expect(response.status).toEqual(200)
    })
  })

  describe('with runs', () => {
    let jobRun: JobRun

    beforeEach(async () => {
      const [node] = await createNuLinkNode(
        db,
        'jobRunsIndexTestNuLinkNode',
      )
      jobRun = await createJobRun(db, node)
    })

    it('returns runs with nulink node names', async () => {
      const response = await request(server).get(
        `/api/v1/job_runs?query=${jobRun.runId}`,
      )
      expect(response.status).toEqual(200)

      const nulinkNode = response.body.included[0]
      expect(nulinkNode.attributes.name).toBeDefined()
      expect(nulinkNode.attributes.accessKey).not.toBeDefined()
      expect(nulinkNode.attributes.salt).not.toBeDefined()
      expect(nulinkNode.attributes.hashedSecret).not.toBeDefined()
    })
  })
})

describe('#show', () => {
  let node: NuLinkNode

  beforeEach(async () => {
    ;[node] = await createNuLinkNode(db, 'jobRunsShowTestNuLinkNode')
  })

  it('returns the job run with task runs', async () => {
    const jobRun = await createJobRun(db, node)
    const response = await request(server).get(`/api/v1/job_runs/${jobRun.id}`)
    expect(response.status).toEqual(200)
    expect(response.body.data.id).toEqual(jobRun.id)
    expect(response.body.data.attributes.runId).toEqual(jobRun.runId)
    expect(response.body.data.relationships.taskRuns.data.length).toEqual(1)
  })

  describe('with out of order task runs', () => {
    let jobRunId: string
    beforeEach(async () => {
      const [nulinkNode] = await createNuLinkNode(
        db,
        'testOutOfOrderTaskRuns',
      )
      const jobRun = new JobRun()
      jobRun.nulinkNodeId = nulinkNode.id
      jobRun.runId = 'OutOfOrderTaskRuns'
      jobRun.jobId = 'xxxxxxxxxxxxxxxxxxxxxxxxxxxxxxxx'
      jobRun.status = 'in_progress'
      jobRun.type = 'runlog'
      jobRun.txHash = 'txA'
      jobRun.requestId = 'requestIdA'
      jobRun.requester = 'requesterA'
      jobRun.createdAt = new Date('2019-04-08T01:00:00.000Z')
      await db.manager.save(jobRun)
      jobRunId = jobRun.id

      const taskRunB = new TaskRun()
      taskRunB.jobRun = jobRun
      taskRunB.index = 1
      taskRunB.status = ''
      taskRunB.type = 'jsonparse'
      await db.manager.save(taskRunB)

      const taskRunA = new TaskRun()
      taskRunA.jobRun = jobRun
      taskRunA.index = 0
      taskRunA.status = 'in_progress'
      taskRunA.type = 'httpget'
      await db.manager.save(taskRunA)
    })

    it('returns ordered task runs', async () => {
      const response = await request(server).get(`/api/v1/job_runs/${jobRunId}`)
      expect(response.status).toEqual(200)
      expect(response.body.data.relationships.taskRuns.data.length).toEqual(2)

      const taskRun1 = response.body.included[1]
      const taskRun2 = response.body.included[2]
      expect(taskRun1.attributes.index).toEqual(0)
      expect(taskRun2.attributes.index).toEqual(1)
    })
  })

  it('returns the job run with only public nulink node information', async () => {
    const jobRun = await createJobRun(db, node)

    const response = await request(server).get(`/api/v1/job_runs/${jobRun.id}`)
    expect(response.status).toEqual(200)

    const clnode = response.body.included[0]
    expect(clnode).toBeDefined()
    expect(clnode.id).toBeDefined()
    expect(clnode.attributes.name).toEqual('jobRunsShowTestNuLinkNode')
    expect(clnode.attributes.accessKey).not.toBeDefined()
    expect(clnode.attributes.hashedSecret).not.toBeDefined()
    expect(clnode.attributes.salt).not.toBeDefined()
  })

  it('returns a 404', async () => {
    const response = await request(server).get('/api/v1/job_runs/1')
    expect(response.status).toEqual(404)
  })
})
