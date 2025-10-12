#!/usr/bin/env python3
"""
Script to refactor i18n function calls to support per-request language.
This adds 'lang string' parameter to all i18n helper functions.
"""

import re
import sys

def refactor_helper_function(content):
    """Add lang parameter to i18n helper function signatures and use TWithLang"""

    # Pattern: func TFoo(params...) string {
    # Replace with: func TFoo(lang string, params...) string {
    pattern = r'func (T[A-Z]\w+)\((.*?)\) string {'

    def replacer(match):
        func_name = match.group(1)
        params = match.group(2).strip()

        # Add lang as first parameter
        if params:
            new_params = f'lang string, {params}'
        else:
            new_params = 'lang string'

        return f'func {func_name}({new_params}) string {{'

    content = re.sub(pattern, replacer, content)

    # Replace i18n.T( with i18n.TWithLang(lang,
    content = re.sub(r'i18n\.T\(', 'i18n.TWithLang(lang, ', content)

    return content

def refactor_handler_function(content):
    """Add lang := h.extractLanguage() to handler functions and update calls"""

    # Find handler functions: func (h *Handlers) FooBar(params...) (*mcp.CallToolResult, error) {
    pattern = r'func \(h \*Handlers\) ([A-Z]\w+)\(params map\[string\]interface\{\}\) \(\*mcp\.CallToolResult, error\) \{'

    def replacer(match):
        func_name = match.group(0)
        # Add language extraction at the start
        return func_name + '\n\tlang := h.extractLanguage()'

    content = re.sub(pattern, replacer, content)

    # Replace i18n helper calls to include lang parameter
    # Pattern: i18n.TFoo( -> i18n.TFoo(lang,
    helper_funcs = [
        'TCommon', 'TParam', 'TValidation', 'TRequiredParam', 'TRequiredArray',
        'TRequiredObject', 'TEntityNotFound', 'TEntityNotFoundByID',
        'TOperationFailed', 'TSuccess', 'TInvalidValue', 'TInvalidStatus',
        'TInvalidPriority', 'TInvalidDate', 'TInvalidFormat', 'TCreateFailed',
        'TUpdateFailed', 'TDeleteFailed', 'TFetchFailed', 'TSaveFailed',
        'TLoadFailed', 'TSearchFailed', 'TMarkdownLabel', 'TListItem',
        'TStatus', 'TPriority', 'TAddFailed', 'TRemoveFailed'
    ]

    for func in helper_funcs:
        # Replace i18n.TFoo( with i18n.TFoo(lang,
        pattern = f'i18n\\.{func}\\('
        replacement = f'i18n.{func}(lang, '
        content = re.sub(pattern, replacement, content)

    # Replace bare i18n.T( calls with i18n.TWithLang(lang,
    content = re.sub(r'(?<!With)i18n\.T\(', 'i18n.TWithLang(lang, ', content)

    return content

if __name__ == '__main__':
    if len(sys.argv) != 3:
        print("Usage: refactor_i18n.py <helpers|handlers> <file>")
        sys.exit(1)

    mode = sys.argv[1]
    filepath = sys.argv[2]

    with open(filepath, 'r') as f:
        content = f.read()

    if mode == 'helpers':
        new_content = refactor_helper_function(content)
    elif mode == 'handlers':
        new_content = refactor_handler_function(content)
    else:
        print(f"Unknown mode: {mode}")
        sys.exit(1)

    # Write to stdout
    print(new_content)
