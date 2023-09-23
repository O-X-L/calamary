from datetime import datetime

# pylint: disable=W0622

project = 'Calamary'
copyright = f'{datetime.now().year}, Superstes'
author = 'Superstes'
extensions = ['sphinx_rtd_theme']
templates_path = ['_templates']
exclude_patterns = []
html_theme = 'sphinx_rtd_theme'
html_theme_options = {}
html_static_path = ['_static']
html_css_files = ['css/main.css']
master_doc = 'index'
display_version = True
sticky_navigation = True
# html_logo = '_static/img/logo.svg'
# html_favicon = 'img/favicon.ico'
