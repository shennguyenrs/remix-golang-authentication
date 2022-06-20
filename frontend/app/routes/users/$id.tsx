import type { LoaderFunction } from '@remix-run/node';
import { useLoaderData } from '@remix-run/react';
import axios from "axios"
import * as libs from "../../../libs";

export const loader: LoaderFunction = async ({ params }) => {
	const {id} = params
	const res = await axios.get(libs.routes.user as string + "/" + id)

	if (res) {
		console.log(res.data)
		return res.data
	} else {
		return { name: "no body" }
	}
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
