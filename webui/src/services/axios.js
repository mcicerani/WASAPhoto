import axios from "axios";

const instance = axios.create({
	baseURL: __API_URL__,
	timeout: 1000 * 5
});

axios.interceptors.request.use(
	config => {
	  const token = localStorage.getItem('token');
	  if (token) {
		config.headers['Authorization'] = token;
	  }
	  return config;
	},
	error => {
	  return Promise.reject(error);
	}
  );
  

export default instance;
