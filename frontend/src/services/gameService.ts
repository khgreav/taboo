import type { AxiosInstance } from 'axios';
import axios from 'axios';

export interface CreateGameParams {
  playerName: string;
}

export interface CreateGameResponse {
  gameId: string;
}

export class GameService {

  private readonly axiosInstance: AxiosInstance;

  /**
   * Constructor
   */
  public constructor() {
    this.axiosInstance = axios.create({
      baseURL: 'http://localhost:8080/api',
      timeout: 30000,
    });
  }

  public async createGame(params: CreateGameParams): Promise<CreateGameResponse> {
    return (await this.axiosInstance.post<CreateGameResponse>('/games', params)).data;
  }
}
