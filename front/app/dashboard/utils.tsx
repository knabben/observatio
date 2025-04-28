// eslint-disable-next-line
export function FilterItems(query: string, items: any[]) {
  return items.filter((i: { name: string; }) =>
    i.name.toLowerCase().includes(query.toLowerCase())).at(0);
}
