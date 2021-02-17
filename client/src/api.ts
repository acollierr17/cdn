import axios from 'axios';
import type { ContextProps } from './contexts/AuthContext';

export interface ImageResult {
  cdn_url: string;
  spaces_url: string;
  spaces_cdn: string;
  file_name: string;
  last_modified: string;
  size: number;
}

export interface ImageResults {
  images: Array<ImageResult>;
  length: number;
}

export interface ImageDeletedResult {
  image_name: string;
  deleted: boolean;
}

export const API_URL = import.meta.env.SNOWPACK_PUBLIC_API_URL;

export const generateToken = async (
  ctx: Partial<ContextProps>,
): Promise<void> => {
  const token = await ctx.user!.getIdToken(true);
  await axios.post(`${API_URL}/token`, { token });
};

export const getAllImages = (
  ctx: Partial<ContextProps>,
  token: string,
): Promise<ImageResults> => {
  return axios
    .get(`${API_URL}/images`, {
      headers: { 'Access-Token': token },
    })
    .then((res) => res.data);
};

export const deleteImage = (
  name: string,
  token: string,
): Promise<ImageDeletedResult> => {
  return axios
    .delete(`${API_URL}/delete/${name}`, {
      headers: { 'Access-Token': token },
    })
    .then((res) => res.data);
};
