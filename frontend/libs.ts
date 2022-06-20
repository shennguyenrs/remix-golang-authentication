import { createCookie } from '@remix-run/node';

export const routes = {
  auth: process.env.BACKEND_AUTH,
  user: process.env.BACKEND_USER,
};

const oneWeek = 1000 * 60 * 60 * 24 * 7;

export const userPref = createCookie('userPref', {
  maxAge: oneWeek,
  expires: new Date(Date.now() + oneWeek),
  path: '/',
  httpOnly: true,
  sameSite: 'strict',
  secure: process.env.NODE_ENV === 'production',
});
