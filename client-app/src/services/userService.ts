import axios from 'axios';

export interface UserRegistrationData {
  name: string;
  email: string;
  password: string;
  role: string;
}

export const createUser = async (userData: UserRegistrationData): Promise<void> => {
  await axios.post('http://localhost:8080/api/users', userData);
}; 