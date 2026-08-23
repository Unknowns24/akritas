"use client";

import React, { useState, useEffect } from "react";
import styles from "./DokploySettingsClient.module.css";
import { Button } from "@/core/ui/primitives/Button";
import { Modal } from "@/core/ui/primitives/Modal";
import { Badge } from "@/core/ui/primitives/Badge";
import { toast } from "sonner";
import { DokployServerForm } from "./components/DokployServerForm";
import { Plus, Trash2, Edit2, CheckCircle, AlertCircle, Clock } from "lucide-react";
import { DokployIcon } from "@/core/ui/icons";
import { DokployServer, listDokployServersService } from "../../services/dokploy/list-dokploy-servers.service";
import { createDokployServerService } from "../../services/dokploy/create-dokploy-server.service";
import { updateDokployServerService } from "../../services/dokploy/update-dokploy-server.service";
import { deleteDokployServerService } from "../../services/dokploy/delete-dokploy-server.service";

export const DokploySettingsClient: React.FC = () => {
  const [servers, setServers] = useState<DokployServer[]>([]);
  const [isModalOpen, setIsModalOpen] = useState(false);
  const [editingServer, setEditingServer] = useState<DokployServer | undefined>();
  const [isLoading, setIsLoading] = useState(false);
  const [error, setError] = useState<string | null>(null);

  useEffect(() => {
    refreshServers();
  }, []);

  const refreshServers = async () => {
    const { data } = await listDokployServersService();
    if (data) setServers(data);
  };

  const handleOpenCreate = () => {
    setEditingServer(undefined);
    setError(null);
    setIsModalOpen(true);
  };

  const handleOpenEdit = (server: DokployServer) => {
    setEditingServer(server);
    setError(null);
    setIsModalOpen(true);
  };

  const handleSubmit = async (formData: { name: string; base_url: string; api_credential?: string }) => {
    setIsLoading(true);
    setError(null);

    let res;
    if (editingServer) {
      res = await updateDokployServerService(editingServer.id, formData);
    } else {
      if (!formData.api_credential) {
        setError("API Credential is required for new servers.");
        setIsLoading(false);
        return;
      }
      res = await createDokployServerService({
        name: formData.name,
        base_url: formData.base_url,
        api_credential: formData.api_credential
      });
    }

    if (res.error) {
      setError(res.error.message || "An error occurred");
      toast.error(res.error.message || "An error occurred");
      setIsLoading(false);
      return;
    }

    await refreshServers();
    toast.success(`Dokploy server ${editingServer ? "updated" : "connected"} successfully`);
    setIsLoading(false);
    setIsModalOpen(false);
  };

  const handleDelete = async (server: DokployServer) => {
    if (server.application_count > 0) {
      const confirmMsg = `WARNING: This server has ${server.application_count} active application(s).\n\nDeleting it may break current incident remediations and references. Are you absolutely sure you want to proceed?`;
      if (!confirm(confirmMsg)) return;
    } else {
      if (!confirm("Are you sure you want to remove this Dokploy server?")) return;
    }

    await deleteDokployServerService(server.id);
    await refreshServers();
    toast.success("Dokploy server removed");
  };

  return (
    <div className={styles.container}>
      <div className={styles.header}>
        <div>
          <h2 className={styles.title}>Dokploy Servers</h2>
          <p className={styles.subtitle}>
            Connect your Dokploy instances to allow Akritas to gather runtime logs and deploy remediations.
          </p>
        </div>
        <Button variant="primary" onClick={handleOpenCreate} className={styles.addBtn}>
          <Plus size={16} /> Connect Server
        </Button>
      </div>

      <div className={styles.list}>
        {servers.length === 0 ? (
          <div className={styles.emptyState}>
            <DokployIcon size={48} className={styles.emptyIcon} />
            <p>No Dokploy servers connected yet.</p>
          </div>
        ) : (
          servers.map((server) => (
            <div key={server.id} className={styles.serverCard}>
              <div className={styles.serverInfo}>
                <div className={styles.avatar}>
                  <DokployIcon size={24} />
                </div>
                <div className={styles.details}>
                  <div className={styles.nameRow}>
                    <span className={styles.name}>{server.name}</span>
                    <Badge variant={server.credential_configured ? "intel" : "neutral"}>
                      {server.credential_configured ? "Credential Configured" : "No Credential"}
                    </Badge>
                  </div>
                  <div className={styles.metaRow}>
                    <span className={styles.metaText}>
                      {server.base_url}
                    </span>
                    <span className={styles.dot}>•</span>
                    <span className={styles.metaText}>
                      {server.application_count} App{server.application_count !== 1 ? 's' : ''}
                    </span>
                    {server.last_synced_at && (
                      <>
                        <span className={styles.dot}>•</span>
                        <span className={styles.metaText}>
                          Synced: {new Date(server.last_synced_at).toLocaleDateString()}
                        </span>
                      </>
                    )}
                  </div>
                </div>
              </div>

              <div className={styles.serverRight}>
                <div className={styles.serverStatus}>
                  {server.connection_status === "connected" && (
                    <span className={styles.statusConnected}><CheckCircle size={14} /> Connected</span>
                  )}
                  {server.connection_status === "authentication_failed" && (
                    <span className={styles.statusFailed}><AlertCircle size={14} /> Auth Failed</span>
                  )}
                  {server.connection_status === "unavailable" && (
                    <span className={styles.statusFailed}><AlertCircle size={14} /> Unavailable</span>
                  )}
                  {server.connection_status === "pending" && (
                    <span className={styles.statusPending}><Clock size={14} /> Pending</span>
                  )}
                </div>

                <div className={styles.serverActions}>
                  <Button variant="ghost" size="sm" onClick={() => handleOpenEdit(server)} className={styles.editBtn}>
                    <Edit2 size={16} />
                  </Button>
                  <Button variant="ghost" size="sm" onClick={() => handleDelete(server)} className={styles.deleteBtn}>
                    <Trash2 size={16} />
                  </Button>
                </div>
              </div>
            </div>
          ))
        )}
      </div>

      <Modal 
        isOpen={isModalOpen} 
        onClose={() => setIsModalOpen(false)}
        title={editingServer ? "Edit Dokploy Server" : "Connect Dokploy Server"}
      >
        {error && <div className={styles.errorAlert}>{error}</div>}
        <DokployServerForm 
          initialData={editingServer}
          onSubmit={handleSubmit}
          onCancel={() => setIsModalOpen(false)}
          isLoading={isLoading}
        />
      </Modal>
    </div>
  );
};
