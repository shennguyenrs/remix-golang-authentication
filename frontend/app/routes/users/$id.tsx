import type { LoaderFunction } from '@remix-run/node';
import type { ErrorBoundaryComponent } from '@remix-run/react/routeModules';
import { useLoaderData } from '@remix-run/react';
import axios from 'axios';
import * as libs from '../../../libs';

export const loader: LoaderFunction = async ({ params, request }) => {
  const { id } = params;
  const cookieHeader = request.headers.get('Cookie');
  const cookie = (await libs.userPref.parse(cookieHeader)) || {};

  try {
    const res = await axios.get((libs.routes.user as string) + '/' + id, {
      headers: {
        Authorization: 'Bearer ' + cookie.token,
      },
    });

    if (res) {
      return res.data;
    } else {
      return { name: 'no body' };
    }
  } catch (err) {
    if (axios.isAxiosError(err) && err.response) {
      throw new Error(err.response.data as string);
    }

    throw new Error('An error occurred while fetching user');
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

export default function UserPage() {
  const { name } = useLoaderData();

  return (
    <div className="rootWrapper">
      <h1>Welcome back, {name}</h1>
      <button className="btnLink-base">Logout</button>
      <button className="btnLink-base">Delete account</button>
    </div>
  );
}
