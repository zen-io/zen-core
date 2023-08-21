
export function mergeMaps(...maps: { [key: string]: string }[]): { [key: string]: string } {
  let result: { [key: string]: string } = {};
  for (let map of maps) {
    for (let key of Object.keys(map)) {
      result[key] = map[key]
    }
  }
  return result;
}
