import { contract, matchers, setup } from '@nulink/test-helpers'
import { assert } from 'chai'
import { ethers } from 'ethers'
import { MedianTestHelperFactory } from '../../ethers/v0.6/MedianTestHelperFactory'

const medianTestHelperFactory = new MedianTestHelperFactory()
const provider = setup.provider()

let defaultAccount: ethers.Wallet
beforeAll(async () => {
  const users = await setup.users(provider)
  defaultAccount = users.roles.defaultAccount
})

describe('Median', () => {
  let median: contract.Instance<MedianTestHelperFactory>

  beforeEach(async () => {
    median = await medianTestHelperFactory.connect(defaultAccount).deploy()
  })

  describe('testing various lists', () => {
    const tests = [
      {
        name: 'ordered ascending',
        responses: [0, 1, 2, 3, 4, 5, 6, 7],
        want: 3,
      },
      {
        name: 'ordered descending',
        responses: [7, 6, 5, 4, 3, 2, 1, 0],
        want: 3,
      },
      {
        name: 'unordered length 1',
        responses: [20],
        want: 20,
      },
      {
        name: 'unordered length 2',
        responses: [20, 0],
        want: 10,
      },
      {
        name: 'unordered length 3',
        responses: [20, 0, 16],
        want: 16,
      },
      {
        name: 'unordered length 4',
        responses: [20, 0, 15, 16],
        want: 15,
      },
      {
        name: 'unordered length 7',
        responses: [1001, 1, 101, 10, 11, 0, 111],
        want: 11,
      },
      {
        name: 'unordered length 9',
        responses: [8, 8, 4, 5, 5, 7, 9, 5, 9],
        want: 7,
      },
      {
        name: 'unordered long',
        responses: [33, 44, 89, 101, 67, 7, 23, 55, 88, 324, 0, 88],
        want: 61, // 67 + 55 / 2
      },
      {
        name: 'unordered longer',
        responses: [
          333121,
          323453,
          337654,
          345363,
          345363,
          333456,
          335477,
          333323,
          332352,
          354648,
          983260,
          333856,
          335468,
          376987,
          333253,
          388867,
          337879,
          333324,
          338678,
        ],
        want: 335477,
      },
      {
        name: 'overflowing numbers',
        responses: [
          ethers.utils.bigNumberify(
            '57896044618658097711785492504343953926634992332820282019728792003956564819967',
          ),
          ethers.utils.bigNumberify(
            '57896044618658097711785492504343953926634992332820282019728792003956564819967',
          ),
        ],
        want: ethers.utils.bigNumberify(
          '57896044618658097711785492504343953926634992332820282019728792003956564819967',
        ),
      },
      {
        name: 'overflowing numbers',
        responses: [
          ethers.utils.bigNumberify(
            '57896044618658097711785492504343953926634992332820282019728792003956564819967',
          ),
          ethers.utils.bigNumberify(
            '57896044618658097711785492504343953926634992332820282019728792003956564819966',
          ),
        ],
        want: ethers.utils.bigNumberify(
          '57896044618658097711785492504343953926634992332820282019728792003956564819966',
        ),
      },
      {
        name: 'really long',
        responses: [
          56,
          2,
          31,
          33,
          55,
          38,
          35,
          12,
          41,
          47,
          21,
          22,
          40,
          39,
          10,
          32,
          49,
          3,
          54,
          45,
          53,
          14,
          20,
          59,
          1,
          30,
          24,
          6,
          5,
          37,
          58,
          51,
          46,
          17,
          29,
          7,
          27,
          9,
          43,
          8,
          34,
          42,
          28,
          23,
          57,
          0,
          11,
          48,
          52,
          50,
          15,
          16,
          26,
          25,
          4,
          36,
          19,
          44,
          18,
          13,
        ],
        want: 29,
      },
    ]

    for (const test of tests) {
      it(test.name, async () => {
        matchers.bigNum(test.want, await median.publicGet(test.responses))
      })
    }
  })

  // long running (minutes) exhaustive test.
  // skipped because very slow, but useful for thorough validation
  xit('permutations', async () => {
    const permutations = (list: number[]) => {
      const result: number[][] = []
      const used: number[] = []

      const permute = (unused: number[]) => {
        if (unused.length == 0) {
          result.push([...used])
          return
        }

        for (let i = 0; i < unused.length; i++) {
          const elem = unused.splice(i, 1)[0]
          used.push(elem)
          permute(unused)
          unused.splice(i, 0, elem)
          used.pop()
        }
      }

      permute(list)
      return result
    }

    {
      const list = [0, 2, 5, 7, 8, 10]
      for (const permuted of permutations(list)) {
        for (let i = 0; i < list.length; i++) {
          for (let j = 0; j < list.length; j++) {
            if (i < j) {
              const foo = await median.publicQuickselectTwo(permuted, i, j)
              matchers.bigNum(list[i], foo[0])
              matchers.bigNum(list[j], foo[1])
            }
          }
        }
      }
    }

    {
      const list = [0, 1, 1, 1, 2]
      for (const permuted of permutations(list)) {
        for (let i = 0; i < list.length; i++) {
          for (let j = 0; j < list.length; j++) {
            if (i < j) {
              const foo = await median.publicQuickselectTwo(permuted, i, j)
              matchers.bigNum(list[i], foo[0])
              matchers.bigNum(list[j], foo[1])
            }
          }
        }
      }
    }
  })

  // Checks the validity of the sorting network in `shortList`
  describe('validate sorting network', () => {
    const net = [
      [0, 1],
      [1, 2],
      [0, 1],
      [3, 4],
      [6, 7],
      [4, 5],
      [7, 8],
      [3, 4],
      [6, 7],
      [0, 3],
      [3, 6],
      [0, 3],
      [1, 4],
      [4, 7],
      [1, 4],
      [5, 8],
      [2, 5],
      [2, 4],
      [4, 6],
      [2, 4],
      [1, 3],
      [2, 3],
      // These last three comparators are commented out in the contract,
      // because they cannot affect the median.
      [5, 8],
      [5, 7],
      [5, 6],
    ]

    // See: https://en.wikipedia.org/wiki/Sorting_network#Zero-one_principle
    xit('zero-one principle', async () => {
      const sortWithNet = (list: number[]) => {
        for (const [i, j] of net) {
          if (list[i] > list[j]) {
            ;[list[i], list[j]] = [list[j], list[i]]
          }
        }
      }

      for (let n = 0; n < (1 << 9) - 1; n++) {
        const list = [
          (n >> 8) & 1,
          (n >> 7) & 1,
          (n >> 6) & 1,
          (n >> 5) & 1,
          (n >> 4) & 1,
          (n >> 3) & 1,
          (n >> 2) & 1,
          (n >> 1) & 1,
          (n >> 0) & 1,
        ]
        const sum = list.reduce((a, b) => a + b, 0)
        sortWithNet(list)
        const sortedSum = list.reduce((a, b) => a + b, 0)
        assert.equal(sortedSum, sum, 'Number of zeros and ones changed')
        list.reduce((switched, i) => {
          assert.isTrue(!switched || i != 0, 'error at n=' + n.toString())
          return i != 0
        }, false)
      }
    })
  })
})
