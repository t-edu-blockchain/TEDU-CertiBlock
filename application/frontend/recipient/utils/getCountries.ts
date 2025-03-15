import { BACKEND_URL } from "./env";

export type Country = { id: number, name: string };

export async function getCountries() {
    const res = await fetch(`${BACKEND_URL}/api/countries`, {
        method: 'GET',
    });

    if (res.status === 200) {
        return await res.json() as Array<Country>;
    }

    return [];
}