import Namable from "./Namable";
import Identifiable from "./Identifiable";

export default interface SearchResult<T extends Namable> {
  type: string;
  results: T[];
}

export interface PrintableResult extends Identifiable, Namable {}
