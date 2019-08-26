configdir for Golang
=====================

Multi platform library of configuration directory for Golang.

This library helps to get regular directories for configuration files or cache files that matches target operationg system's convention.

It assumes the following folders are standard paths of each environment:

.. list-table::
   :header-rows: 1

   - * 
     * Windows:
     * Linux/BSDs:
     * MacOSX:
   - * System level configuration folder
     * ``%PROGRAMDATA%`` (``C:\\ProgramData``)
     * ``${XDG_CONFIG_DIRS}`` (``/etc/xdg``)
     * ``/Library/Application Support``
   - * User level configuration folder
     * ``%APPDATA%`` (``C:\\Users\\<User>\\AppData\\Roaming``)
     * ``${XDG_CONFIG_HOME}`` (``${HOME}/.config``)
     * ``${HOME}/Library/Application Support``
   - * User wide cache folder
     * ``%LOCALAPPDATA%`` ``(C:\\Users\\<User>\\AppData\\Local)``
     * ``${XDG_CACHE_HOME}`` (``${HOME}/.cache``)
     * ``${HOME}/Library/Caches``

Examples
------------

Getting Configuration
~~~~~~~~~~~~~~~~~~~~~~~~

``configdir.ConfigDir.QueryFolderContainsFile()`` searches files in the following order:

* Local path (if you add the path via LocalPath parameter)
* User level configuration folder(e.g. ``$HOME/.config/<vendor-name>/<application-name>/setting.json`` in Linux)
* System level configuration folder(e.g. ``/etc/xdg/<vendor-name>/<application-name>/setting.json`` in Linux)

``configdir.Config`` provides some convenient methods(``ReadFile``, ``WriteFile`` and so on).

.. code-block:: go

   var config Config

   configDirs := configdir.New("vendor-name", "application-name")
   // optional: local path has the highest priority
   configDirs.LocalPath, _ = filepath.Abs(".")
   folder := configDirs.QueryFolderContainsFile("setting.json")
   if folder != nil {
       data, _ := folder.ReadFile("setting.json")
       json.Unmarshal(data, &config)
   } else {
       config = DefaultConfig
   }

Write Configuration
~~~~~~~~~~~~~~~~~~~~~~

When storing configuration, get configuration folder by using ``configdir.ConfigDir.QueryFolders()`` method.

.. code-block:: go

   configDirs := configdir.New("vendor-name", "application-name")

   var config Config
   data, _ := json.Marshal(&config)

   // Stores to local folder
   folders := configDirs.QueryFolders(configdir.Local)
   folders[0].WriteFile("setting.json", data)

   // Stores to user folder
   folders = configDirs.QueryFolders(configdir.Global)
   folders[0].WriteFile("setting.json", data)

   // Stores to system folder
   folders = configDirs.QueryFolders(configdir.System)
   folders[0].WriteFile("setting.json", data)

Getting Cache Folder
~~~~~~~~~~~~~~~~~~~~~~

It is similar to the above example, but returns cache folder.

.. code-block:: go

   configDirs := configdir.New("vendor-name", "application-name")
   cache := configDirs.QueryCacheFolder()

   resp, err := http.Get("http://examples.com/sdk.zip")
   if err != nil {
       log.Fatal(err)
   }
   defer resp.Body.Close()
   body, err := ioutil.ReadAll(resp.Body)

   cache.WriteFile("sdk.zip", body)

Document
------------

https://godoc.org/github.com/shibukawa/configdir

License
------------

MIT

