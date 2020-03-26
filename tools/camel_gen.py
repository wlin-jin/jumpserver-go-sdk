def snake_to_camel(word):
        import re
        return ''.join(x.capitalize() or '_' for x in word.split('_'))

def check_field_type(value):
	v_type = type(v).__name__
	if v_type == 'str':
		return "string"
	elif v_type == 'int':
		return "int"
	elif  v_type == 'bool':
		return "bool"
	else:
		return "string"

tpl = '''{camel_name} {field_type} `json:"{json_name}"`'''
for k, v in data.items():
	info = {
		"camel_name": snake_to_camel(k),
		"field_type": check_field_type(v),
		"json_name": k
	}
	print(tpl.format(**info))
