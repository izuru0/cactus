import type { Express, Request, Response } from "express";

import {
  Logger,
  Checks,
  LogLevelDesc,
  LoggerProvider,
  IAsyncProvider,
} from "@hyperledger/cactus-common";
import {
  IEndpointAuthzOptions,
  IExpressRequestHandler,
  IWebServiceEndpoint,
} from "@hyperledger/cactus-core-api";
import { registerWebServiceEndpoint } from "@hyperledger/cactus-core";

import OAS from "../../json/openapi.json";
import { PluginKeychainVault } from "../plugin-keychain-vault";

export interface ISetKeychainEntryEndpointV1Options {
  logLevel?: LogLevelDesc;
  plugin: PluginKeychainVault;
}

export class SetKeychainEntryEndpointV1 implements IWebServiceEndpoint {
  public static readonly CLASS_NAME = "SetKeychainEntryEndpointV1";

  private readonly log: Logger;
  private readonly plugin: PluginKeychainVault;

  public get className(): string {
    return SetKeychainEntryEndpointV1.CLASS_NAME;
  }

  constructor(public readonly options: ISetKeychainEntryEndpointV1Options) {
    const fnTag = `${this.className}#constructor()`;
    Checks.truthy(options, `${fnTag} arg options`);
    Checks.truthy(options.plugin, `${fnTag} arg options.plugin`);
    Checks.truthy(
      options.plugin instanceof PluginKeychainVault,
      `${fnTag} arg options.plugin instanceof PluginKeychainVault`,
    );

    this.plugin = options.plugin;

    const level = this.options.logLevel || "INFO";
    const label = this.className;
    this.log = LoggerProvider.getOrCreate({ level, label });
    this.log.debug(`Instantiated ${this.className} OK`);
  }

  public get oasPath(): (typeof OAS.paths)["/api/v1/plugins/@hyperledger/cactus-plugin-keychain-vault/set-keychain-entry"] {
    return OAS.paths[
      "/api/v1/plugins/@hyperledger/cactus-plugin-keychain-vault/set-keychain-entry"
    ];
  }

  getAuthorizationOptionsProvider(): IAsyncProvider<IEndpointAuthzOptions> {
    // TODO: make this an injectable dependency in the constructor
    return {
      get: async () => ({
        isProtected: true,
        requiredRoles: [],
      }),
    };
  }

  public async registerExpress(
    expressApp: Express,
  ): Promise<IWebServiceEndpoint> {
    await registerWebServiceEndpoint(expressApp, this);
    return this;
  }

  public getVerbLowerCase(): string {
    return this.oasPath.post["x-hyperledger-cactus"].http.verbLowerCase;
  }

  public getPath(): string {
    return this.oasPath.post["x-hyperledger-cactus"].http.path;
  }

  public getOperationId(): string {
    return this.oasPath.post.operationId;
  }
  public getExpressRequestHandler(): IExpressRequestHandler {
    return this.handleRequest.bind(this);
  }

  public async handleRequest(req: Request, res: Response): Promise<void> {
    const tag = `${this.getVerbLowerCase().toUpperCase()} ${this.getPath()}`;
    try {
      this.log.debug(`${tag} %o`, req.body);
      const { key, value } = req.body;
      const resBody = await this.plugin.set(key, value);
      res.status(200);
      res.json(resBody);
    } catch (ex) {
      this.log.debug(`${tag} Failed to serve request:`, ex);
      res.status(500);
      res.json({ error: ex.stack });
    }
  }
}
