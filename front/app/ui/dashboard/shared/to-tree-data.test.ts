import {toTreeData} from './to-tree-data';

describe('toTreeData', () => {
  it('converts a scalar leaf', () => {
    const tree = toTreeData({name: 'c1'});
    expect(tree).toEqual([{label: 'name: c1', value: 'name'}]);
  });

  it('converts a nested object into an expandable node with unique child paths', () => {
    const tree = toTreeData({spec: {paused: false}});
    expect(tree).toHaveLength(1);
    expect(tree[0].label).toBe('spec');
    expect(tree[0].value).toBe('spec');
    expect(tree[0].children).toEqual([{label: 'paused: false', value: 'spec.paused'}]);
  });

  it('converts an array into indexed nodes', () => {
    const tree = toTreeData({items: ['a', 'b']});
    expect(tree[0].children).toEqual([
      {label: '[0]: a', value: 'items.[0]'},
      {label: '[1]: b', value: 'items.[1]'},
    ]);
  });

  it('renders an empty object as a leaf, not a misleadingly-expandable node', () => {
    const tree = toTreeData({metadata: {}});
    expect(tree).toEqual([{label: 'metadata: {}', value: 'metadata'}]);
  });

  it('renders an empty array as a leaf', () => {
    const tree = toTreeData({conditions: []});
    expect(tree).toEqual([{label: 'conditions: []', value: 'conditions'}]);
  });

  it('renders null/undefined values as a placeholder, not a crash', () => {
    const tree = toTreeData({a: null, b: undefined});
    expect(tree).toEqual([
      {label: 'a: —', value: 'a'},
      {label: 'b: —', value: 'b'},
    ]);
  });

  it('returns an empty array for a top-level scalar or null input', () => {
    expect(toTreeData('just a string')).toEqual([]);
    expect(toTreeData(null)).toEqual([]);
  });

  it('handles deep nesting with unique paths at every level', () => {
    const tree = toTreeData({status: {conditions: [{type: 'Ready', status: 'True'}]}});
    const conditionsNode = tree[0].children?.[0];
    expect(conditionsNode?.value).toBe('status.conditions');
    const firstCondition = conditionsNode?.children?.[0];
    expect(firstCondition?.value).toBe('status.conditions.[0]');
    expect(firstCondition?.children).toEqual([
      {label: 'type: Ready', value: 'status.conditions.[0].type'},
      {label: 'status: True', value: 'status.conditions.[0].status'},
    ]);
  });
});
