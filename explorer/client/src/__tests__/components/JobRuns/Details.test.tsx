import { mount } from 'enzyme'
import React from 'react'
import { partialAsFull } from '@nulink/ts-test-helpers'
import { JobRun, NuLinkNode } from 'explorer/models'
import Details from '../../../components/JobRuns/Details'

describe('components/JobRuns/Details', () => {
  it('hides error when not present', () => {
    const nulinkNode = partialAsFull<NuLinkNode>({})
    const jobRun = partialAsFull<JobRun>({ nulinkNode })

    const wrapper = mount(<Details jobRun={jobRun} etherscanHost="" />)

    expect(wrapper.text()).not.toContain('Error')
  })

  it('displays error when present', () => {
    const nulinkNode = partialAsFull<NuLinkNode>({})
    const jobRun = partialAsFull<JobRun>({ error: 'Failure!', nulinkNode })

    const wrapper = mount(<Details jobRun={jobRun} etherscanHost="" />)

    expect(wrapper.text()).toContain('Error')
    expect(wrapper.text()).toContain('Failure!')
  })
})
