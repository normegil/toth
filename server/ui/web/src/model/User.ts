import Role from "./Role";

export default class User {
  id: string;
  mail: string;
  name: string;
  role: Role;

  constructor(id: string, mail: string, name: string, role: Role) {
    this.id = id;
    this.mail = mail;
    this.name = name;
    this.role = role;
  }
}
