package mcwig

var lspServerInitJson = `{
  "processId": null,
  "rootPath": "",
  "clientInfo": {
    "name": "mcwig",
    "version": "001"
  },
  "rootUri": "",
  "capabilities": {
    "general": {
      "positionEncodings": [
        "utf-32",
        "utf-16"
      ]
    },
    "workspace": {
      "workspaceEdit": {
        "documentChanges": true,
        "resourceOperations": [
          "create",
          "rename",
          "delete"
        ]
      },
      "applyEdit": true,
      "symbol": {
        "symbolKind": {
          "valueSet": [
            1,
            2,
            3,
            4,
            5,
            6,
            7,
            8,
            9,
            10,
            11,
            12,
            13,
            14,
            15,
            16,
            17,
            18,
            19,
            20,
            21,
            22,
            23,
            24,
            25,
            26
          ]
        }
      },
      "executeCommand": {
        "dynamicRegistration": false
      },
      "didChangeWatchedFiles": {
        "dynamicRegistration": true
      },
      "workspaceFolders": true,
      "configuration": true,
      "diagnostics": {
        "refreshSupport": false
      },
      "fileOperations": {
        "didCreate": false,
        "willCreate": false,
        "didRename": true,
        "willRename": true,
        "didDelete": false,
        "willDelete": false
      }
    },
    "textDocument": {
      "declaration": {
        "dynamicRegistration": true,
        "linkSupport": true
      },
      "references": {
        "dynamicRegistration": true
      },
      "implementation": {
        "dynamicRegistration": true,
        "linkSupport": true
      },
      "typeDefinition": {
        "dynamicRegistration": true,
        "linkSupport": true
      },
      "synchronization": {
        "willSave": true,
        "didSave": true,
        "willSaveWaitUntil": true
      },
      "documentSymbol": {
        "symbolKind": {
          "valueSet": [
            1,
            2,
            3,
            4,
            5,
            6,
            7,
            8,
            9,
            10,
            11,
            12,
            13,
            14,
            15,
            16,
            17,
            18,
            19,
            20,
            21,
            22,
            23,
            24,
            25,
            26
          ]
        },
        "hierarchicalDocumentSymbolSupport": true
      },
      "formatting": {
        "dynamicRegistration": true
      },
      "rangeFormatting": {
        "dynamicRegistration": true
      },
      "onTypeFormatting": {
        "dynamicRegistration": true
      },
      "rename": {
        "dynamicRegistration": true,
        "prepareSupport": true
      },
      "codeAction": {
        "dynamicRegistration": true,
        "isPreferredSupport": true,
        "codeActionLiteralSupport": {
          "codeActionKind": {
            "valueSet": [
              "",
              "quickfix",
              "refactor",
              "refactor.extract",
              "refactor.inline",
              "refactor.rewrite",
              "source",
              "source.organizeImports"
            ]
          }
        },
        "resolveSupport": {
          "properties": [
            "edit",
            "command"
          ]
        },
        "dataSupport": true
      },
      "completion": {
        "completionItem": {
          "snippetSupport": true,
          "documentationFormat": [
            "markdown",
            "plaintext"
          ],
          "resolveAdditionalTextEditsSupport": true,
          "insertReplaceSupport": true,
          "deprecatedSupport": true,
          "resolveSupport": {
            "properties": [
              "documentation",
              "detail",
              "additionalTextEdits",
              "command",
              "insertTextFormat",
              "insertTextMode"
            ]
          },
          "insertTextModeSupport": {
            "valueSet": [
              1,
              2
            ]
          }
        },
        "contextSupport": true,
        "dynamicRegistration": true
      },
      "signatureHelp": {
        "signatureInformation": {
          "parameterInformation": {
            "labelOffsetSupport": true
          }
        },
        "dynamicRegistration": true
      },
      "documentLink": {
        "dynamicRegistration": true,
        "tooltipSupport": true
      },
      "hover": {
        "contentFormat": [
          "markdown",
          "plaintext"
        ],
        "dynamicRegistration": true
      },
      "selectionRange": {
        "dynamicRegistration": true
      },
      "callHierarchy": {
        "dynamicRegistration": false
      },
      "typeHierarchy": {
        "dynamicRegistration": true
      },
      "publishDiagnostics": {
        "relatedInformation": true,
        "tagSupport": {
          "valueSet": [
            1,
            2
          ]
        },
        "versionSupport": true
      },
      "diagnostic": {
        "dynamicRegistration": false,
        "relatedDocumentSupport": false
      },
      "linkedEditingRange": {
        "dynamicRegistration": true
      }
    },
    "window": {
      "workDoneProgress": true,
      "showDocument": {
        "support": true
      }
    }
  },
  "initializationOptions": null,
  "workDoneToken": "1"
}`


