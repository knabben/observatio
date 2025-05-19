// eslint-disable-next-line
export function FilterItems(query: string, items: any[]) {
  return items.filter((i: { metadata: {name: string} }) =>
    i.metadata?.name.toLowerCase().includes(query.toLowerCase())).at(0);
}
