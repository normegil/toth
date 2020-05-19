import Role from "./Role";

export default class User {
  id: string;
  name: string;
  role: Role;

  constructor(id: string, name: string, role: Role) {
    this.id = id;
    this.name = name;
    this.role = role;
  }
}
