import 'vue3-toastify/dist/index.css';
import { type ToastContainerOptions } from 'vue3-toastify';

export const ToastOptions: ToastContainerOptions = {
  position: 'top-right',
  autoClose: 5000,
  closeButton: true,
  pauseOnHover: true,
  pauseOnFocusLoss: true,
  closeOnClick: true,
  theme: 'colored',
  limit: 1,
}

export { default } from 'vue3-toastify';
