import { Form } from '@remix-run/react';
import type { ErrorBoundaryComponent } from '@remix-run/react/routeModules';
import type { ActionFunction } from '@remix-run/node';
import { redirect } from '@remix-run/node';
import axios from 'axios';
import * as libs from '../../libs';

export const action: ActionFunction = async ({ request }) => {
  const form = await request.formData();
  const name = form.get('username');
  const email = form.get('email');
  const password = form.get('password');
  const confirmPassword = form.get('confirmPassword');
  const cookieHeader = request.headers.get('Cookie');

  if (password !== confirmPassword) {
    throw new Error('Passwords do not match');
  }

  try {
    const res = await axios.post((libs.routes.auth as string) + '/register', {
      name,
      email,
      password,
    });

    if (res) {
      const { token, userId } = res.data;
      const cookie = (await libs.userPref.parse(cookieHeader)) || {};
      cookie.token = token;

      return redirect('/users/' + userId, {
        headers: {
          'Set-Cookie': await libs.userPref.serialize(cookie),
        },
      });
    }
  } catch (err) {
    if (axios.isAxiosError(err) && err.response) {
      throw new Error(err.response.data as string);
    }

    throw new Error('An error occurred while registering new user');
  }
};

export const ErrorBoundary: ErrorBoundaryComponent = ({ error }) => {
  return (
    <div className="rootWrapper">
      <p>{error.message}</p>
      <a className="btnLink-base text-white bg-blue-500 mt-5" href="/">
        Return home
      </a>
    </div>
  );
};

export default function Register() {
  return (
    <div className="rootWrapper">
      <h1>Register page</h1>
      <Form className="formCon" method="post">
        <label>Username</label>
        <input type="text" name="username" />
        <label>Email</label>
        <input type="email" name="email" />
        <label>Password</label>
        <input type="password" name="password" />
        <label>Confirm Password</label>
        <input type="password" name="confirmPassword" />
        <button
          className="btnLink-base text-white bg-blue-500 font-bold"
          type="submit"
        >
          Register
        </button>
      </Form>
      <div>
        <a className="btnLink-base text-white bg-red-300" href="/">
          Go back to homepage
        </a>
      </div>
    </div>
  );
}
