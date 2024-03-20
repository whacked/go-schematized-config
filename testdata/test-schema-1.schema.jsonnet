{
  title: 'example-properties-schema',
  description: 'example JSON Schema for demonstration and testing in documentation using nbdev',
  type: 'object',
  required: [
    'string_value_with_enum',
    'MY_INTEGER_VALUE',
    'A_NUMERIC_VALUE',
  ],
  properties: {
    // string
    SOME_STRING_VALUE_funnyCaSe345: {
      type: 'string',
    },
    string_value_with_enum: {
      description: 'this one is in the <required> list!',
      type: 'string',
      enum: [
        'it',
        'can',
        'only',
        'be',
        'one',
        'of',
        'these',
      ],
    },
    _____A_STRING_VALUE____with_default__: {
      description: 'values with a default get hydrated using the default if not present in input',
      type: 'string',
      default: 'underscores_and spaces',
    },

    // integer
    MY_INTEGER_VALUE: {
      type: 'integer',
      description: 'not used for validation, but your benefit',
    },

    // number
    A_NUMERIC_VALUE: {
      type: 'number',
      description: 'continuous and real and reasonable',
      minimum: 22,
      maximum: 33333.4,
    },

    // boolean
    true_or_false__but_also_nothing: {
      type: 'boolean',
    },
  },
}
